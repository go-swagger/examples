// SPDX-FileCopyrightText: Copyright 2015-2026 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package restapi

import (
	stderrors "errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"

	"github.com/go-swagger/examples/file-server/restapi/operations"
	"github.com/go-swagger/examples/file-server/restapi/operations/uploads"
)

type uploadHTTPResult struct {
	statusCode int
	err        error
}

type streamingHandlerProgress struct {
	firstChunk string
	fields     url.Values
	files      []runtime.MultipartFileInfo
}

type streamingHandlerResult struct {
	fields url.Values
	files  []runtime.MultipartFileInfo
	err    error
}

func TestStreamingUploadReachesHandlerBeforeRequestBodyCompletes(t *testing.T) {
	const (
		firstChunk  = "first"
		secondChunk = "second"
	)

	swaggerSpec, err := loads.Embedded(SwaggerJSON, FlatSwaggerJSON)
	require.NoError(t, err)

	handlerRead := make(chan streamingHandlerProgress, 1)
	handlerDone := make(chan streamingHandlerResult, 1)

	api := operations.NewFileUploadAPI(swaggerSpec)
	api.MultipartformConsumer = runtime.DiscardConsumer
	api.UploadsUploadFileHandler = uploads.UploadFileHandlerFunc(func(params uploads.UploadFileParams) middleware.Responder {
		defer func() {
			_ = params.MultipartForm.Close()
		}()

		file, handlerErr := params.MultipartForm.NextFile()
		if handlerErr == nil {
			first := make([]byte, len(firstChunk))
			_, handlerErr = io.ReadFull(file, first)
			if handlerErr == nil {
				handlerRead <- streamingHandlerProgress{
					firstChunk: string(first),
					fields:     params.MultipartForm.Fields(),
					files:      params.MultipartForm.Files(),
				}
				_, handlerErr = io.Copy(io.Discard, file)
			}
		}
		if handlerErr == nil {
			_, handlerErr = params.MultipartForm.NextFile()
			if stderrors.Is(handlerErr, io.EOF) {
				handlerErr = nil
			}
		}
		if handlerErr == nil {
			handlerErr = params.MultipartForm.Drain()
		}

		handlerDone <- streamingHandlerResult{
			fields: params.MultipartForm.Fields(),
			files:  params.MultipartForm.Files(),
			err:    handlerErr,
		}
		if handlerErr != nil {
			return middleware.Error(http.StatusInternalServerError, handlerErr)
		}

		return uploads.NewUploadFileOK()
	})

	server := httptest.NewServer(api.Serve(nil))
	defer server.Close()

	bodyReader, bodyWriter := io.Pipe()
	defer func() {
		_ = bodyReader.Close()
	}()
	defer func() {
		_ = bodyWriter.Close()
	}()

	multipartWriter := multipart.NewWriter(bodyWriter)
	request, err := http.NewRequestWithContext(t.Context(), http.MethodPost, server.URL+"/upload", bodyReader)
	require.NoError(t, err)
	request.Header.Set(runtime.HeaderContentType, multipartWriter.FormDataContentType())

	responseDone := make(chan uploadHTTPResult, 1)
	go func() {
		response, requestErr := server.Client().Do(request)
		if requestErr != nil {
			responseDone <- uploadHTTPResult{err: requestErr}

			return
		}
		defer func() {
			_ = response.Body.Close()
		}()

		responseDone <- uploadHTTPResult{statusCode: response.StatusCode}
	}()

	require.NoError(t, multipartWriter.WriteField("before", "one"))
	part, err := multipartWriter.CreateFormFile("file", "payload.bin")
	require.NoError(t, err)
	_, err = io.WriteString(part, firstChunk)
	require.NoError(t, err)

	select {
	case progress := <-handlerRead:
		assert.EqualT(t, firstChunk, progress.firstChunk)
		assert.EqualT(t, "one", progress.fields.Get("before"))
		assert.Empty(t, progress.fields.Get("after"))
		require.Len(t, progress.files, 1)
		assert.EqualT(t, "file", progress.files[0].FieldName)
		assert.EqualT(t, "payload.bin", progress.files[0].Filename)
	case result := <-responseDone:
		require.NoError(t, result.err)
		t.Fatal("request completed before the multipart body was resumed")
	case <-time.After(time.Second):
		t.Fatal("handler did not receive the first file chunk while the request body was still open")
	}

	_, err = io.WriteString(part, secondChunk)
	require.NoError(t, err)
	require.NoError(t, multipartWriter.WriteField("after", "two"))
	require.NoError(t, multipartWriter.Close())
	require.NoError(t, bodyWriter.Close())

	var result uploadHTTPResult
	select {
	case result = <-responseDone:
	case <-time.After(time.Second):
		t.Fatal("request did not complete after the multipart body was closed")
	}
	require.NoError(t, result.err)
	assert.EqualT(t, http.StatusOK, result.statusCode)

	select {
	case handlerResult := <-handlerDone:
		require.NoError(t, handlerResult.err)
		assert.EqualT(t, "one", handlerResult.fields.Get("before"))
		assert.EqualT(t, "two", handlerResult.fields.Get("after"))
		require.Len(t, handlerResult.files, 1)
		assert.EqualT(t, "payload.bin", handlerResult.files[0].Filename)
	case <-time.After(time.Second):
		t.Fatal("upload handler did not complete")
	}
}
