---
title: "go-swagger examples"
type: home
description: 'Runnable examples and tutorials for go-swagger spec-first code generation'
weight: 1
---

A curated collection of **runnable examples** for
[`go-swagger`](https://github.com/go-swagger/go-swagger) — generating servers,
clients and CLIs from an OpenAPI 2.0 (Swagger) spec.

Every example here is committed to the
[go-swagger/examples](https://github.com/go-swagger/examples) repository and kept
in sync with the latest go-swagger release by automated regeneration.

### Status

{{% button href="https://github.com/go-swagger/examples/fork" hint="fork me on github" style=primary icon=code-fork %}}Fork me{{% /button %}}
Actively maintained. Regenerated weekly against `swagger@master`.

### Which site do I want?

These examples are all **spec-first**: you have an OpenAPI spec and want
`swagger generate` to produce typed code. The sibling sites cover the other two
approaches — pick by what you start from:

| I start from… | I want… | Go here |
|---------------|---------|---------|
| an **OpenAPI spec** | generate a typed server / client / CLI | **this site** |
| **Go interfaces**, no codegen | hand-wire an untyped client or server | [go-openapi/runtime](https://go-openapi.github.io/runtime/) |
| **Go code** | produce a spec *from* the code (code-first) | [go-openapi/codescan](https://go-openapi.github.io/codescan/) |

### New to go-swagger?

Install the toolchain and read the command reference on go-swagger's own site:

```cmd
go install github.com/go-swagger/go-swagger/cmd/swagger@latest
```

→ [go-swagger.io](https://goswagger.io/go-swagger/) for install, the `generate`
command families, and project-layout reference. This site assumes you have
`swagger` on your `PATH` and focuses on **what to build with it**.

### Where to go next

{{< cards >}}
{{% card title="Guides" %}}
The example catalog, grouped by concern — servers, clients & CLI, authentication,
streaming, and codegen customization. One page per example.

→ [guides](./guides/)
{{% /card %}}

{{% card title="Tutorials" %}}
Sequential, end-to-end walkthroughs. Start with the todo-list tutorial to build a
server and client from scratch.

→ [tutorials](./tutorials/)
{{% /card %}}

{{% card title="Project" %}}
Repository README, licensing, contributing guidelines and how the examples stay in
sync with go-swagger.

→ [project](./project/)
{{% /card %}}
{{< /cards >}}

## Licensing

`SPDX-FileCopyrightText: Copyright 2025 go-swagger maintainers`

These examples ship under the [Apache-2.0 license](./project/LICENSE.md).

## Contributing

Issues and pull requests welcome. See [project/](./project/) for guidelines.

---

{{< children type="card" description="true" >}}
