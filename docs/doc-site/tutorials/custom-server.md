---
title: "Custom server tutorial"
weight: 3
description: "Embed a generated core in a custom server"
---

The [todo-list tutorial](../todo-list/) generated a *whole* server — `main.go`
and all — and you edited the `configure_*.go` file it left for you. Sometimes you
want the opposite balance: keep go-swagger's generated **core** (the models,
router, and typed operations) but own the `main` yourself, so the CLI is a thin
hand-written layer that wires configuration and handlers around that core.

That's what `--exclude-main` is for. This tutorial builds a tiny *greeter* server
that way.

{{% notice info %}}
You'll need the `swagger` CLI on your `PATH` — see
[goswagger.io](https://goswagger.io/go-swagger/install/). The finished code lives
under [`tutorials/custom-server/`](https://github.com/go-swagger/examples/tree/master/tutorials/custom-server).
{{% /notice %}}

## Step 1 — the spec

The greeter is deliberately minimal: one `GET /hello` that takes an optional
`name` query parameter and returns a plain-text greeting.

```yaml
swagger: '2.0'
info:
  version: 1.0.0
  title: Greeting Server
paths:
  /hello:
    get:
      produces:
        - text/plain
      parameters:
        - name: name
          required: false
          type: string
          in: query
          description: defaults to World if not given
      operationId: getGreeting
      responses:
        200:
          description: returns a greeting
          schema:
            type: string
            description: contains the actual greeting as plain text
```

## Step 2 — generate the core only

Generate the server into a `gen/` sub-tree with `--exclude-main`, so go-swagger
emits everything *except* a `main.go`:

```sh
rm -rf gen && mkdir gen
swagger generate server --exclude-main -A greeter -t gen -f ./swagger/swagger.yml
```

You get `gen/restapi/` — the embedded spec, the `NewServer` constructor, the
router, and `operations/` with the typed `GreeterAPI`, its `GetGreetingParams`,
and the `GetGreetingOK` responder. What you *don't* get is a `cmd/` entry point.
That's yours to write.

## Step 3 — write your own `main`

Your `main` does what the generated `main.go` would have — load the embedded spec,
construct the API, and hand it to `NewServer` — but it's plain code you control:

{{< code file="tutorials/custom-server/cmd/greeter/main.go" lang="go" region="wiring" >}}

Because you own this file, you can add your own flags, config loading, dependency
injection, logging, or lifecycle management around this core — none of it is
generated, none of it gets overwritten on regeneration.

## Step 4 — attach the handler

The generated `GreeterAPI` exposes one handler field per operation. Assign your
implementation to `GetGreetingHandler` before serving — this is the same handler
you'd otherwise place in a generated `configure_*.go`, but here it lives in your
`main`:

{{< code file="tutorials/custom-server/cmd/greeter/main.go" lang="go" region="handler" >}}

`conv.Value` dereferences the optional `*string` parameter, defaulting to `World`,
and `NewGetGreetingOK().WithPayload(...)` returns the typed `200` responder the
generated code defined.

## Step 5 — run it

```shellsession
$ go run ./cmd/greeter/main.go --port 3000
```

Then exercise it (here with [httpie](https://httpie.org)):

```shellsession
$ http get :3000/hello                  # Hello, World!
$ http get :3000/hello name==Swagger    # Hello, Swagger!
```

## Regenerating safely

The whole point of the split is that regeneration only ever touches `gen/`. When
the spec changes, rerun the Step 2 command — `gen/` is rewritten, your `cmd/`
`main` is untouched. Keep the two apart (generated core under `gen/`, your code
outside it) and the two never collide.

## Related

- [Todo list tutorial](../todo-list/) — the opposite balance: generate the whole server and edit `configure_*.go`.
- [Custom middleware guide](../../guides/customizing-codegen/custom-middleware/) — extend a *fully* generated server via its hook points, no `--exclude-main` needed.
- [Generation flags guide](../../guides/customizing-codegen/generation-flags/) — other flags that reshape the generated `main`.
