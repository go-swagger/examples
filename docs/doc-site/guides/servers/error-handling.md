---
title: "Custom error handling"
weight: 3
description: "Customizing error responses in a generated server"
---

By default a generated handler returns a `middleware.Responder` and you build
error responses yourself (`NewAddOneDefault(500).WithPayload(...)`). The
**todo-list-errors** example shows the alternative: generate with
`--return-errors` so handlers may return a plain `error`, and install a single
custom error handler that shapes every error response in one place.

{{% notice tip %}}
Source: [`todo-list-errors/`](https://github.com/go-swagger/examples/tree/master/todo-list-errors).
Generated with `swagger generate server -A TodoList -f ./swagger.yml --return-errors`
(the `//go:generate` line in `configure_todo_list.go` uses the long-form flags).
{{% /notice %}}

## Handlers that return an error

With `--return-errors`, the handler signature becomes
`func(params) (middleware.Responder, error)`. A handler can now short-circuit by
returning an `error` instead of constructing a response:

{{< code file="todo-list-errors/restapi/configure_todo_list.go" lang="go" region="handler" >}}

Here `errAlreadyExists` is a sentinel the handler returns directly — no response
plumbing at the call site.

## The centralized error handler

Returned errors flow through `api.ServeError`, which you can override. This
example wires it to a `catcher` that recognizes the sentinel (with
`errors.Is`), logs it, then defers to the runtime's default `ServeError` for the
actual HTTP response:

{{< code file="todo-list-errors/restapi/configure_todo_list.go" lang="go" region="catcher" >}}

Installing it is one line in `configureAPI`:

```go
api.ServeError = catcher
```

Because every operation's error funnels through the same hook, you get one place
to classify errors, attach correlation IDs, translate domain errors to status
codes, or emit metrics — without repeating that logic in each handler.

## When to use it

Use returned errors when your handlers naturally surface Go `error` values (a
data layer, a validation step) and you'd rather map them to responses centrally
than build a `*Default` responder at every return. Stick with the default
responder style when each handler already knows the exact response it wants.

## Related

- [Todo list server](../todo-list/) — the default responder style.
- [Strict server](../strict-server/) — compiler-enforced responders.
- Hand-wiring error handling without codegen? See the
  [go-openapi/runtime](https://go-openapi.github.io/runtime/) examples.
