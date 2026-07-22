---
title: "Todo list server"
weight: 1
description: "Canonical full server: unix, http and https listeners"
---

The **todo-list** example is the canonical go-swagger server: a small CRUD API
generated from a spec, wired to a trivial in-memory store. It's the best place to
see what `swagger generate server` produces and which files you're expected to
edit.

Prefer a step-by-step build? Start with the [todo-list tutorial](../../../tutorials/todo-list/).
This page is the reference tour of the finished example.

{{% notice tip %}}
Source: [`todo-list/`](https://github.com/go-swagger/examples/tree/master/todo-list).
Regenerate with `go run ./hack/tools regen` (or the `//go:generate` directive in
`restapi/configure_todo_list.go`).
{{% /notice %}}

## The spec

Two definitions drive everything: an `item` (with a required, `minLength: 1`
`description` and a read-only `id`) and a generic `error`.

{{< code file="todo-list/swagger.yml" lang="yaml" region="definitions" >}}

The `readOnly: true` on `id` means the server assigns it — clients that send one
have it ignored on create.

## Generated models

Each definition becomes a Go struct. Note how the spec constraints map onto the
type: `description` is required and `minLength: 1`, so it's generated as a
non-pointer-safe `*string` carrying validation, while the read-only `id` is a
plain `int64`.

{{< code file="todo-list/models/item.go" lang="go" lines="17-30" >}}

The generator also emits the validation methods that enforce the spec's
constraints at runtime — here the `description` field's `required` and
`minLength: 1` rules, wired straight from the spec. You never hand-write this:

{{< code file="todo-list/models/item.go" lang="go" lines="46-57" >}}

## Wiring the handlers

`restapi/configure_todo_list.go` is the one file you edit — it's marked *safe to
edit* and survives regeneration. The generated scaffold leaves each handler
returning `501 Not Implemented`; you replace those with real logic. Here's the
implemented version from the tutorial's `server-complete`, backing the store with
a map:

{{< code file="tutorials/todo-list/server-complete/restapi/configure_todo_list.go" lang="go" region="handlers" >}}

Each response constructor (`NewAddOneCreated`, `NewDestroyOneNoContent`, …) is
generated from a response you declared in the spec, so the compiler keeps your
handlers honest against the contract.

## Running it

The generated server listens on a Unix socket, HTTP and HTTPS by default. For a
quick local test, enable just the HTTP listener on a fixed port:

```shellsession
$ go run ./todo-list/cmd/todo-list-server --scheme=http --port=8765
serving todo list at http://127.0.0.1:8765
```

```shellsession
$ curl -i localhost:8765 \
    -H 'Content-Type: application/io.goswagger.examples.todo-list.v1+json' \
    -d '{"description":"go shopping"}'
HTTP/1.1 201 Created
...
{"description":"go shopping","id":1}
```

See the [todo-list README](https://github.com/go-swagger/examples/blob/master/todo-list/README.md)
for the full listener matrix (unix/http/https) and TLS options.

## Variations

- [Strict server](../strict-server/) — the strict responder interface.
- [Custom error handling](../error-handling/) — shaping error responses.
- [Todo list tutorial](../../../tutorials/todo-list/) — build this from scratch.
