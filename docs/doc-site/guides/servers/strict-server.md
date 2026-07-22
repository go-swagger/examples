---
title: "Strict server"
weight: 2
description: "Server generated with the strict responder interface"
---

The **strict server** is generated with `--strict-responders`. Instead of every
handler returning a generic `middleware.Responder` — which lets you hand back
*any* response, including ones the spec never declared — each operation gets its
own **responder interface**. The compiler then enforces that a handler can only
return responses declared for that operation.

{{% notice tip %}}
Source: [`todo-list-strict/`](https://github.com/go-swagger/examples/tree/master/todo-list-strict).
Generated with `swagger generate server -A todo-list -f ./swagger.yml --strict-responders --regenerate-configureapi`.
{{% /notice %}}

## The generated responder interface

For each operation, the generator emits a marker interface embedding
`middleware.Responder`. Only the response types declared for `addOne` implement
`AddOneResponder`, so nothing else can be returned:

{{< code file="todo-list-strict/restapi/operations/todos/add_one_responses.go" lang="go" lines="127-130" >}}

Every generated response for the operation — `AddOneCreated`, `AddOneDefault`,
`AddOneNotImplemented` — carries a no-op `AddOneResponder()` method, which is
what admits it to the interface. A `FindDefault` value, for instance, simply
won't compile inside an `addOne` handler.

## The handler signature

Compare with the [default todo-list server](../todo-list/), where the handler
returns `middleware.Responder`. Here the return type is the operation-specific
`todos.AddOneResponder`:

{{< code file="todo-list-strict/restapi/configure_simple_to_do_list_api.go" lang="go" region="strict-handler" >}}

The scaffold guards each assignment with `if … == nil`, so you can wire a handler
from elsewhere and leave the rest returning `NotImplemented`. You replace the
body with real logic, returning one of the operation's typed responders.

## When to use it

Reach for strict responders when you want the type system to guarantee your
handlers stay in sync with the contract — you can't accidentally return a
response shape the spec doesn't describe. The cost is a little more generated
surface (one interface per operation) and slightly more verbose returns.

## Related

- [Todo list server](../todo-list/) — the default (non-strict) responder style.
- [Custom error handling](../error-handling/) — returning errors from handlers.
