---
title: "Auto-configure"
weight: 4
description: "Auto-wire handler implementations to the generated API"
---

Normally you edit `configure_*.go` by hand to attach each handler. The
**auto-configure** example takes a different route: generate with
`--implementation-package`, point it at a package you write, and the generator
emits the wiring for you. Every operation is routed to a method on your
implementation — no hand-editing of a configure file.

{{% notice tip %}}
Source: [`auto-configure/`](https://github.com/go-swagger/examples/tree/master/auto-configure).
Generated with `swagger generate server --name AToDoListApplication --spec ./swagger.yml --implementation-package github.com/go-swagger/examples/auto-configure/implementation --principal any`.
{{% /notice %}}

## The generated contract

Instead of a `configure_*.go`, the generator produces
`auto_configure_*.go`. It declares the `Handler` interface your package must
satisfy, and binds a package-level `Impl` to your constructor via
`implementation.New()`:

{{< code file="auto-configure/restapi/auto_configure_a_to_do_list_application.go" lang="go" lines="24-55" >}}

Each generated handler then simply delegates to that `Impl`:

{{< code file="auto-configure/restapi/auto_configure_a_to_do_list_application.go" lang="go" lines="76-79" >}}

## Your implementation package

You write a package that implements `Handler`. This is ordinary hand-written code
— nothing generated — so it's the natural home for your business logic. Here's
the `AddOne` method backing the in-memory store:

{{< code file="auto-configure/implementation/todos_impl.go" lang="go" region="add-one" >}}

The example splits the interface across small types — `TodosHandlerImpl`,
`ConfigureImpl`, `AuthImpl` — composed by a single `HandlerImpl` that `New()`
returns. That keeps handlers, server configuration, and authentication in
separate files while still satisfying the one generated `Handler` interface.

## When to use it

Auto-configure shines when you regenerate often and don't want a hand-edited
`configure_*.go` in the loop: your implementation lives entirely in a package you
own, and regeneration only ever touches the generated wiring. It's also a clean
way to keep the transport-facing glue separate from your domain code.

## Related

- [Todo list server](../todo-list/) — the conventional hand-wired `configure_*.go`.
- [Custom error handling](../error-handling/) — centralized error responses.
