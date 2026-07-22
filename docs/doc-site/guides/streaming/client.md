---
title: "Streaming client"
weight: 2
description: "Consume a stream from a generated client"
---

A generated client normally reads the *entire* response body and unmarshals it
into a typed payload. To consume a stream you override that behavior — swap the
consumer, and either buffer the whole thing or read it chunk-by-chunk. Two
examples show both approaches.

{{% notice tip %}}
Sources: [`stream-server/elapsed_client.go`](https://github.com/go-swagger/examples/blob/master/stream-server/elapsed_client.go)
(non-blocking, pairs with the [streaming server](../server/)) and
[`stream-client/jigsaw.go`](https://github.com/go-swagger/examples/blob/master/stream-client/jigsaw.go)
(blocking, against an external server).
{{% /notice %}}

## Why the default doesn't stream

The generated client hands the response body to a **consumer**, which does an
`io.Copy` into the destination you pass. With the default `JSONConsumer` that
copy blocks until the body is complete and then unmarshals — exactly what you
*don't* want for a stream. The fix is to install a `ByteStreamConsumer` for the
response mime, so the bytes flow through untouched.

## Non-blocking: consume chunks as they arrive

Override the consumer, then pass an `io.Pipe` writer as the destination. The
client's `io.Copy` writes into the pipe while a goroutine reads the other end —
so bytes are processed the moment they arrive:

{{< code file="stream-server/elapsed_client.go" lang="go" region="consumer" >}}

A `bufio.Scanner` on the read end splits the stream on newlines and unmarshals
each line independently. `cancel()` (via `defer`) tears down the request if the
reader stops early:

{{< code file="stream-server/elapsed_client.go" lang="go" region="scan" >}}

The request runs on the main goroutine, writing into the pipe the scanner drains;
a context timeout bounds how long it will keep the connection open:

{{< code file="stream-server/elapsed_client.go" lang="go" region="request" >}}

## Blocking: buffer the whole stream

If you don't need incremental processing, the simplest path is a destination that
knows how to accept `text/plain`. Give a buffer an `UnmarshalText` method and the
default consumer will fill it:

{{< code file="stream-client/jigsaw.go" lang="go" region="buffer" >}}

Then pass it straight to the operation and read it once the call returns:

{{< code file="stream-client/jigsaw.go" lang="go" region="blocking" >}}

The `jigsaw.go` example also has a non-blocking variant that installs
`transport.Consumers[runtime.TextMime] = runtime.ByteStreamConsumer()` — the same
technique as above, for a `text/plain` stream instead of newline-delimited JSON.

## Choosing an approach

- **Blocking + `UnmarshalText`** — least code; fine when you can wait for the full
  response.
- **Non-blocking + `ByteStreamConsumer` + `io.Pipe`** — process items as they
  arrive, cancel early, bound with a context. Use this for long-lived or unbounded
  streams.

## Related

- [Streaming server](../server/) — the countdown server `elapsed_client.go` consumes.
