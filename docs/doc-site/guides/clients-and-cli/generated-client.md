---
title: "Generated client SDK"
weight: 1
description: "Classic vs stratoscale client flavors"
---

From the same spec that drives a server, `swagger generate client` produces a
typed Go SDK: one method per operation, with generated parameter and response
types. The **tutorials/client** example generates that SDK in two flavors from
one spec — the **classic** go-swagger client and the **stratoscale** contributed
template — so you can compare the ergonomics.

{{% notice tip %}}
Source: [`tutorials/client/`](https://github.com/go-swagger/examples/tree/master/tutorials/client).
Two client packages are generated from one `swagger.yml`:

- classic — `swagger generate client -A TodoList --spec swagger.yml --client-package classic-client`
- stratoscale — the same, plus `--template stratoscale` (reusing `--existing-models`).
{{% /notice %}}

## Classic client

The classic client generates a `ClientService` interface. Each method takes the
operation's params plus an explicit `runtime.ClientAuthInfoWriter` and variadic
`ClientOption`s; a parallel `…Context` variant threads a `context.Context`:

{{< code file="tutorials/client/classic_client/todos/todos_client.go" lang="go" lines="101-103" >}}

## Stratoscale client

The stratoscale template generates a leaner `API` interface: **context-first**,
auth folded into the transport, no options parameter. It also emits a
`//go:generate mockery` directive so the interface is trivially mockable in tests:

{{< code file="tutorials/client/stratoscale_client/todos/todos_client.go" lang="go" lines="18-19" >}}

Pick classic for the full go-swagger surface (per-call auth, per-call options);
pick stratoscale for a compact, context-first, easily-mocked client.

## Multiple success responses

This example deliberately exercises a tricky spec shape: `addOne` declares **two**
success responses (`201 Created` and `204 No Content`). Both flavors reflect that
in the return signature — the method hands back a pointer for *each* possible
success, and exactly one is non-nil:

```go
created, noContent, err := c.AddOne(ctx, params)
switch {
case err != nil:
    // transport or error response
case created != nil:
    // 201 — use created.Payload
case noContent != nil:
    // 204 — nothing to read
}
```

The same mechanism covers operations with **no default response**: without a
`default`, an undeclared status code surfaces as a generic error rather than a
typed payload, because the generated response reader only knows the codes the spec
listed.

## Related

- [CLI client](../cli-client/) — a cobra command-line tool wrapping a generated client.
- [Custom templates](../../customizing-codegen/custom-templates/) — how the stratoscale flavor is produced.
- Hand-wiring a client without codegen? See the
  [go-openapi/runtime](https://go-openapi.github.io/runtime/) examples.
