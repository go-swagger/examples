---
title: "Repository README"
weight: 1
description: "The go-swagger/examples repository overview"
---

This site documents the [`go-swagger/examples`](https://github.com/go-swagger/examples)
repository — a collection of runnable, committed examples for
[`go-swagger`](https://github.com/go-swagger/go-swagger), the **spec-first** code
generator for OpenAPI 2.0 (Swagger).

## What's here

Every example is a real Go project generated from an OpenAPI 2.0 spec: servers,
typed client SDKs, and CLIs. The generated code is **committed** and kept in sync
with go-swagger `master` by automated [regeneration](../regeneration/), so what you
read on this site matches what the current generator emits.

Browse the material two ways:

- **[Guides](../../guides/)** — reference recipes you dip into (servers,
  clients & CLI, authentication, streaming, customizing codegen).
- **[Tutorials](../../tutorials/)** — sequential, end-to-end walkthroughs that build
  something from scratch.

## Where this fits — three sibling sites

go-swagger and go-openapi split their example material across three sites by
*workflow*. This one is the **spec-first** corner: you write an OpenAPI spec and
generate typed Go from it.

| Site | Workflow | You start from |
|------|----------|----------------|
| **this site** | spec-first codegen | an OpenAPI 2.0 spec → generated server/client/CLI |
| [go-openapi/runtime](https://go-openapi.github.io/runtime/) | untyped / hand-wired | the runtime API, no codegen |
| [go-openapi/codescan](https://go-openapi.github.io/codescan/) | code-first | Go code → generated spec |

When a topic straddles two workflows, the page links across.

## Getting started

```sh
git clone https://github.com/go-swagger/examples
```

You'll need the `swagger` CLI on your `PATH` to regenerate or follow the
tutorials — see the [installation instructions](https://goswagger.io/go-swagger/install/).
Then head to the [todo-list tutorial](../../tutorials/todo-list/) to build a server
and client from one spec.

## Status & releasing

The examples track go-swagger code generation for servers and clients. The
repository is deliberately left **unreleased**: it follows the generator on
`go-swagger/go-swagger@master` rather than tagging versions of its own.

## Licensing

This software ships under the [Apache-2.0](../license/) license.

## Other documentation

- [Contributing guidelines](../contributing/)
- [Regeneration](../regeneration/) — how the examples stay in sync
- [All-time contributors](https://github.com/go-swagger/examples/blob/master/CONTRIBUTORS.md)
