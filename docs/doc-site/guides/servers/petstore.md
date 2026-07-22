---
title: "Petstore"
weight: 7
description: "The classic Swagger petstore, generated"
---

The **Swagger Petstore** is the canonical OpenAPI 2.0 sample. The `generated/`
example is that spec run through `swagger generate server` — a complete server
scaffold for a multi-resource API, useful as a reference for what a "full" spec
produces.

{{% notice tip %}}
Source: [`generated/`](https://github.com/go-swagger/examples/tree/master/generated).
Generated with `swagger generate server --name Petstore --spec ./swagger-petstore.json --principal any`.
{{% /notice %}}

## Three resource groups

The petstore spec organizes operations under three tags, each becoming its own
package of generated handlers:

- **pet** — `addPet`, `updatePet`, `findPetsByStatus`, `findPetsByTags`,
  `getPetById`, `deletePet`, `uploadFile`
- **store** — `getInventory`, `placeOrder`, `getOrderById`, `deleteOrder`
- **user** — `createUser`, `getUserByName`, `loginUser`, `logoutUser`, …

## A generated model

The `Pet` model shows the usual spec-to-Go mapping — required fields as pointers,
nested models by reference, arrays as slices — alongside a generated `Validate`
method (not shown) that enforces the spec's constraints:

{{< code file="generated/models/pet.go" lang="go" lines="20-42" >}}

## Two security schemes

Petstore mixes an api-key header with an OAuth2 flow. The generator emits a
distinct authenticator hook for each — note the OAuth2 hook receives the required
`scopes` so you can authorize per operation:

{{< code file="generated/restapi/configure_petstore.go" lang="go" region="auth" >}}

See the [OAuth2 guide](../../authentication/oauth2-access-code/) for a fully
implemented flow.

## Typed vs untyped

This repo also ships the **same petstore API hand-wired without codegen**, using
the go-openapi runtime directly, under
[`2.0/petstore`](https://github.com/go-swagger/examples/tree/master/2.0/petstore).
That untyped style — building the API from runtime primitives rather than
generated code — is the subject of the
[go-openapi/runtime](https://go-openapi.github.io/runtime/) site. Compare the two
to see exactly what `swagger generate` buys you.

## Related

- [Todo list server](../todo-list/) — a smaller server to start from.
- [Task tracker](../task-tracker/) — another full spec, generated end to end.
