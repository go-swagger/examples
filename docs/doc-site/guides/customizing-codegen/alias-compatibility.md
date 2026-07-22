---
title: "Alias compatibility"
weight: 5
description: "Type-alias compatibility in generated code"
---

{{% notice info %}}
This example runs the *other* direction — **code → spec** (`swagger generate
spec`), not spec → code. That code-first workflow is the subject of the
[go-openapi/codescan](https://go-openapi.github.io/codescan/) site; it lives here
only because the repo hosts the example. The rest of this site is spec-first
codegen.
{{% /notice %}}

The **alias-compatibility** example shows how a Go **type alias** is reflected
when you generate a spec *from* Go code, and how the `--transparent-aliases` flag
controls it.

{{% notice tip %}}
Source: [`alias-compatibility/`](https://github.com/go-swagger/examples/tree/master/alias-compatibility).
{{% /notice %}}

## The aliases

`UserID` is a true Go alias (`=`) of `Identifier`, not a distinct named type:

{{< code file="alias-compatibility/api.go" lang="go" region="aliases" >}}

## What the flag does

When `swagger generate spec` walks this code, the alias can be treated two ways:

- **Default (post-[#3227](https://github.com/go-swagger/go-swagger/issues/3227))** —
  `UserID` appears as its own definition, and `User.id` references
  `#/definitions/UserID`.
- **`--transparent-aliases`** — `UserID` is *not* emitted; `User.id` references
  `#/definitions/Identifier` directly (the pre-#3227 behavior).

See the difference by generating both and diffing:

```shellsession
$ swagger generate spec -m -o without-flag.json
$ swagger generate spec -m --transparent-aliases -o with-flag.json
$ diff <(jq . without-flag.json) <(jq . with-flag.json)
```

## Related

- [go-openapi/codescan](https://go-openapi.github.io/codescan/) — the code-first (code → spec) workflow this example belongs to.
- [External types](../external-types/) — the spec-first counterpart: binding a schema to an existing Go type.
