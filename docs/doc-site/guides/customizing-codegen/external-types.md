---
title: "External types"
weight: 2
description: "Bind schemas to externally defined Go types"
---

By default every schema in your spec becomes a generated Go struct. Sometimes you
want a schema to map onto a type you already have — a hand-written type, one from
another package, or a shared domain type. The **external-types** example shows how
`x-go-type` binds a schema to an externally defined Go type instead of generating
one.

{{% notice tip %}}
Source: [`external-types/`](https://github.com/go-swagger/examples/tree/master/external-types).
See also the go-swagger [external types reference](https://goswagger.io/go-swagger/use/models/schemas/#external-types).
{{% /notice %}}

## The `x-go-type` extension

Attach `x-go-type` to a schema to name the Go type it should use, and where to
import it from. Here a property is bound to `MyAlternateInteger` from the `fred`
package instead of getting a generated type:

{{< code file="external-types/example-external-types.yaml" lang="yaml" region="x-go-type" >}}

The `import.package` (and optional `alias`) tell the generator which import to add.
The generated code references your type directly — no definition is emitted for it.

## The generated result

A definition bound to an external type collapses to exactly that type, with the
external package imported (and its name mangled to avoid collisions). This
`MyExtCollection` is a slice of an external `go-ext` type:

{{< code file="external-types/models/my_ext_collection.go" lang="go" lines="16-19" >}}

Because the external type is expected to satisfy the runtime's `Validatable`
interface, the generated `Validate` still calls into it per item — so your type
participates in validation like any generated model.

## What it covers

The example exercises the full range of external-type use cases:

- an external type as its own definition, or nested inside an object/slice/map/tuple;
- types pulled from the default models package or from an arbitrary import path;
- embedding an external type to add the `Validatable` interface;
- annotation *hints* to resolve nullable/struct-vs-interface questions and to skip
  validation of an external type.

{{% notice note %}}
The example spec adds an `additionalItems` clause to demonstrate tuples, which
makes it not strictly valid against the Swagger 2.0 meta-schema — intentional, to
show the tuple binding.
{{% /notice %}}

## Related

- [Custom templates](../custom-templates/) — reshape the generated code itself.
- [Generation flags](../generation-flags/) — control model/target packages.
