---
title: "Custom templates"
weight: 1
description: "Generate with contributed (stratoscale) templates"
---

go-swagger renders code from Go templates, and you can swap in your own. The
**contributed-templates** example uses the built-in `--template stratoscale`
option — a community template set that produces a different, interface-first
shape optimized for testability.

{{% notice tip %}}
Source: [`contributed-templates/stratoscale/`](https://github.com/go-swagger/examples/tree/master/contributed-templates/stratoscale).
Generated with `swagger generate server -A Petstore --template stratoscale` (and the
matching `swagger generate client --template stratoscale`).
{{% /notice %}}

## What the template changes

Instead of the default's per-operation handler *func fields*, the stratoscale
template groups operations into **interfaces** — one per tag — that your code
implements, and emits `//go:generate mockery` directives so those interfaces are
trivially mockable in tests:

{{< code file="contributed-templates/stratoscale/restapi/configure_petstore.go" lang="go" lines="26-42" >}}

Every handler is `context`-first (`ctx context.Context, params …`), matching the
[stratoscale client flavor](../../clients-and-cli/generated-client/). You provide
one implementation per interface (`PetAPI`, `StoreAPI`, …) rather than assigning
individual handler funcs.

## When to use it

Reach for a custom template set when the default output doesn't fit your codebase
conventions — here, interface-based handlers plus generated mocks. Because it's a
whole template family, it changes server *and* client output consistently.

To go further, `--template-dir` points the generator at your own template
directory, letting you override any individual template go-swagger ships.

## Related

- [Generated client SDK](../../clients-and-cli/generated-client/) — the stratoscale client interface, side by side with the classic one.
- [Generation flags](../generation-flags/) — other flags that reshape the output.
