# File upload server

This example demonstrates how to build a simple file upload endpoint
with swagger and go-swagger.

## Try it

1. Build the server

```
cd restapi/cmd/file-upload-server
go build

./file-upload-server --port 8000
2021/01/17 18:54:09 Serving file upload at http://127.0.0.1:8000
```

2. Run the client

From another terminal:

```
go run upload_file.go swagger.yml
```

Logs on the server:
```
2021/01/17 18:54:15 received file name: swagger.yml
2021/01/17 18:54:15 received file size: 512
2021/01/17 18:54:15 copied bytes 512
2021/01/17 18:54:15 file uploaded copied as upload427417421/uploaded_file_0.dat
```

The file has been copied in a temporary folder `cmd/file-upload-server/upload*/`


## Specification

We use the swagger type `file` in a multipart form, like so:

```yaml
paths:
  /upload:
    post:
      consumes:
      - multipart/form-data
      parameters:
      - name: file
        in: formData
        type: file
```

## Server side

The handler receives a `*runtime.StreamedFile`, which reads the payload directly
from the HTTP request body. The filename and MIME headers are available before
the payload is consumed:

```go
	log.Printf("received file name: %s", params.File.Filename)
	log.Printf("received content type: %s", params.File.Header.Get(runtime.HeaderContentType))
```

The complete file size is not known before the stream is consumed.

## Experimental server-side streaming

This branch demonstrates the intended server binding for a file parameter with
`x-go-server-streaming: true`.

The generated binder is adapted manually until go-swagger supports the
extension. It uses `runtime.MultipartFormStream` from go-openapi/runtime#507 and
exposes the file payload before the complete multipart request has arrived.
The handler owns the multipart stream and must either:

- call `Drain` after consuming the file to process the remaining parts and close
  the request body; or
- call `Close` to abort multipart processing.

The request size is limited outside the binder with `http.MaxBytesHandler`.
The runtime stream's own body limit is disabled so that `*http.MaxBytesError`
from the outer middleware is propagated through file reads and draining.

## Client side

The local file is handled as a `runtime.NamedReadCloser` (that is, a `io.ReadCloser` plus the `Name() string` method).
A regular `os.File` satisfies this.

The file can be passed directly to the client method, like so:

```go
	params := uploads.NewUploadFileParams().WithFile(reader)

	_, err := uploader.Uploads.UploadFile(params)
```
