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

For `x-go-server-streaming: true`, generated binding only constructs a
`*runtime.MultipartFormStream` and passes ownership to the handler. It does not
consume parts or try to populate generated file and form fields.

The handler traverses the multipart body sequentially:

```go
	for {
		file, err := params.MultipartForm.NextFile()
		if errors.Is(err, io.EOF) {
			break
		}
		// Validate file.FieldName and consume file here.
	}
```

Each `runtime.StreamedFile` reads directly from the HTTP request body. The
filename and MIME headers are available before the payload is consumed, but the
complete file size is not known in advance.

## Server-side streaming

The server binding is generated from `x-go-server-streaming: true`. The
generated binder creates a `*runtime.MultipartFormStream` without reading
multipart parts ahead of the handler. Required fields, accepted file field
names, multiplicity and other application-specific rules are validated by the
handler while traversing the stream.

`MultipartFormStream.Fields()` and `MultipartFormStream.Files()` return
snapshots of ordinary fields and file metadata discovered so far. They do not
read ahead: trailing fields become visible only after the active file is
consumed or closed and the stream advances.

The handler owns the multipart stream and must either:

- call `Drain()` to process all remaining parts and close the request body; or
- call `Close()` to abort multipart processing.

## Client side

The local file is handled as a `runtime.NamedReadCloser` (that is, a `io.ReadCloser` plus the `Name() string` method).
A regular `os.File` satisfies this.

The file can be passed directly to the client method, like so:

```go
	params := uploads.NewUploadFileParams().WithFile(reader)

	_, err := uploader.Uploads.UploadFile(params)
```
