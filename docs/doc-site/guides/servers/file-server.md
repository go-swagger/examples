---
title: "File server"
weight: 5
description: "File upload and download endpoints"
---

The **file-server** example demonstrates a file-upload endpoint: how the spec's
`type: file` maps onto the generated server and client, and how the runtime
surfaces the uploaded file to your handler.

{{% notice tip %}}
Source: [`file-server/`](https://github.com/go-swagger/examples/tree/master/file-server).
Build the server under `restapi/cmd/file-upload-server`, then run the client with
`go run upload_file.go swagger.yml`.
{{% /notice %}}

## The spec

An upload is a `multipart/form-data` operation with a `formData` parameter of
`type: file`:

{{< code file="file-server/swagger.yml" lang="yaml" region="upload-path" >}}

## Server side

The generated handler receives the file as an `io.ReadCloser` on
`params.File`. At runtime it's a `*runtime.File`, so a type assertion gives you
the multipart header — filename and size — before you stream the body to disk:

{{< code file="file-server/restapi/configure_file_upload.go" lang="go" region="upload-handler" >}}

Note the `defer params.File.Close()` and the `io.Copy` into a fresh file — the
handler owns the stream and is responsible for draining and closing it.

## Client side

On the client, a file argument is a `runtime.NamedReadCloser` — an
`io.ReadCloser` that also reports a `Name()`. A plain `*os.File` satisfies it, so
you open the file and pass it straight to the generated parameter builder:

{{< code file="file-server/upload_file.go" lang="go" region="client-upload" >}}

The generated `UploadFile` client method handles the multipart encoding; you only
supply the reader.

## Related

- Streaming request/response bodies instead of a one-shot upload? See
  [Streaming](../../streaming/).
- Hand-wiring multipart without codegen? See the
  [go-openapi/runtime](https://go-openapi.github.io/runtime/) examples.
