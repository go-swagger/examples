---
title: "Streaming server"
weight: 1
description: "Stream newline-delimited JSON bodies from a generated server"
---

Swagger 2.0 has no first-class notion of a streaming response, but you can still
generate a server that streams. The **stream-server** example is a countdown API:
`GET /elapse/{length}` emits one newline-delimited JSON object per second until it
reaches zero.

{{% notice tip %}}
Source: [`stream-server/`](https://github.com/go-swagger/examples/tree/master/stream-server).
Generated with `swagger generate server --spec ./swagger.yml`.
{{% /notice %}}

## Declaring a streaming response

The trick is the response schema: `type: string, format: binary`. That's the
closest Swagger 2.0 gets to "this endpoint streams bytes", and it makes the
generator produce a response the handler writes to directly rather than a typed
payload it serializes for you:

{{< code file="stream-server/swagger.yml" lang="yaml" region="streaming-response" >}}

## Writing the stream from the handler

Instead of returning a generated responder, the handler returns a
`middleware.ResponderFunc` — a closure with raw access to the
`http.ResponseWriter`. It grabs the `http.Flusher` and writes through a small
wrapper so each write is pushed to the client immediately:

{{< code file="stream-server/restapi/configure_countdown.go" lang="go" region="handler" >}}

The `flushWriter` is what turns a normal write into a streamed chunk — it flushes
after every write, so the client sees each line as it's produced rather than at
the end:

{{< code file="stream-server/restapi/configure_countdown.go" lang="go" region="flush-writer" >}}

## Producing the chunks

The business logic just encodes one `Mark` per iteration into the writer, with a
one-second pause between them. Because the writer flushes on every `Encode`, each
`{"remains":N}` line reaches the client in real time:

{{< code file="stream-server/biz/count.go" lang="go" region="producer" >}}

## Trying it

```shellsession
$ go run ./stream-server/cmd/countdown-server --port=8000
$ curl -N http://127.0.0.1:8000/elapse/5
{"remains":5}
{"remains":4}
{"remains":3}
{"remains":2}
{"remains":1}
{"remains":0}
```

The response uses `Transfer-Encoding: chunked`; each line arrives a second apart.
A length of `11` returns `403` (a contrived error to show non-streaming responses
still work normally).

## Related

- [Streaming client](../client/) — consuming this stream from a generated client.
- Hand-wiring a streaming server without codegen? See the
  [go-openapi/runtime](https://go-openapi.github.io/runtime/) examples.
