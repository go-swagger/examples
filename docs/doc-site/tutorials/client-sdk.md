---
title: "Client SDK tutorial"
weight: 2
description: "Generate and use a typed SDK client"
---

The [todo-list tutorial](../todo-list/) built a **server** from a spec. This one
takes the *same* kind of spec and generates a typed **client SDK** — one Go method
per operation, with generated parameter and response types — then walks through
actually calling an API with it.

You'll generate the SDK in two flavors from one spec (the **classic** go-swagger
client and the leaner **stratoscale** contributed template) and see how the
generated signature copes with two spec shapes that trip people up: an operation
with **multiple success responses**, and one with **no default response**.

{{% notice info %}}
You'll need the `swagger` CLI on your `PATH` — see
[goswagger.io](https://goswagger.io/go-swagger/install/). The finished code lives
under [`tutorials/client/`](https://github.com/go-swagger/examples/tree/master/tutorials/client)
in the examples repository. For a side-by-side *reference* comparison of the two
flavors, see the [generated client SDK guide](../../guides/clients-and-cli/generated-client/);
this page is the hands-on walkthrough.
{{% /notice %}}

## Step 1 — the spec

Start from a todo-list `swagger.yml`. The only detail that shapes the client here
is the security scheme: the API is protected by an API-key header, so every
request the client sends must carry a `x-todolist-token`:

{{< code file="tutorials/client/swagger.yml" lang="yaml" region="security" >}}

That `security` requirement is what forces an *auth writer* into the calls below.

## Step 2 — generate the classic client

Point `swagger generate client` at the spec:

```sh
swagger generate client -A TodoList --spec swagger.yml --client-package classic_client
```

This writes a `classic_client/` package: a top-level `TodoList` client whose
fields group the operations by tag (`Todos`, `Experimental`), plus generated
`…Params` and `…Responses` types under each tag package.

## Step 3 — call the API

Instantiate the client over a transport, then call an operation. Because the spec
requires an API key, you pass a `runtime.ClientAuthInfoWriter` built from the
transport's `APIKeyAuth` helper — the classic client takes it **per call**:

```go
import (
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag/conv"

	client "github.com/go-swagger/examples/tutorials/client/classic_client"
	"github.com/go-swagger/examples/tutorials/client/classic_client/todos"
	"github.com/go-swagger/examples/tutorials/client/models"
)

// point the transport at the running server
transport := httptransport.New("localhost:8080", "/", []string{"http"})
c := client.New(transport, strfmt.Default)

// the API-key writer that satisfies the `key` security scheme
auth := httptransport.APIKeyAuth("x-todolist-token", "header", "my-secret-token")

params := todos.NewAddOneParams().WithBody(&models.Item{
	Description: conv.Pointer("write the client tutorial"),
})

created, noContent, err := c.Todos.AddOne(params, auth)
```

Notice the call returns **three** values, not two — that's the next step.

## Step 4 — the stratoscale flavor

Regenerate the same spec with the stratoscale template for a leaner, context-first
client (auth folded into construction, no per-call options, and a
`//go:generate mockery` directive for easy mocking):

```sh
swagger generate client -A TodoList --spec swagger.yml \
  --template stratoscale --existing-models ... --client-package stratoscale_client
```

Usage folds the auth writer into the constructor and threads a `context.Context`
through each call instead of an auth argument:

```go
import (
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"

	client "github.com/go-swagger/examples/tutorials/client/stratoscale_client"
	"github.com/go-swagger/examples/tutorials/client/stratoscale_client/todos"
)

c := client.New(client.Config{
	URL:      mustParse("http://localhost:8080/"),
	AuthInfo: httptransport.APIKeyAuth("x-todolist-token", "header", "my-secret-token"),
})

created, noContent, err := c.Todos.AddOne(ctx, todos.NewAddOneParams().WithBody(item))
```

Same operation, same three return values — only the ergonomics differ. Pick
classic for the full go-swagger surface (per-call auth and options); pick
stratoscale for a compact, mockable client.

## Step 5 — multiple success responses

Why three return values? Because `addOne` declares **two** success responses in
the spec — `201 Created` *and* `204 No Content`:

```yaml
    post:
      operationId: addOne
      responses:
        '201':
          description: Created
          schema:
            $ref: "#/definitions/item"
        '204':
          description: Already there
```

The generated method reflects that by handing back a pointer for *each* possible
success; exactly one is non-nil. Switch on them:

```go
created, noContent, err := c.Todos.AddOne(params, auth)
switch {
case err != nil:
	// transport failure, or a typed error response
case created != nil:
	// 201 — the new item is in created.Payload
	log.Printf("created #%d", created.Payload.ID)
case noContent != nil:
	// 204 — the item already existed; nothing to read
	log.Print("already there")
}
```

The **no default response** case (the `experimental` operations declare `401`/`405`
but no `default`) works by the same mechanism in reverse: with no `default`, any
status code the spec didn't list can't map to a typed payload, so it surfaces as a
generic error from the response reader rather than a `*…Default` value.

## Related

- [Generated client SDK guide](../../guides/clients-and-cli/generated-client/) — the reference comparison of both flavors.
- [CLI client guide](../../guides/clients-and-cli/cli-client/) — wrap a generated client in a cobra command-line tool.
- [Todo list tutorial](../todo-list/) — the server side of the same spec.
