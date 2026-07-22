---
title: "Composed auth"
weight: 2
description: "Composing multiple security requirements"
---

Real APIs rarely have a single security scheme. The **composed-auth** example
mixes four — basic auth, an API key (by header *or* query), and scoped JWT
tokens — and composes them per operation with **AND** and **OR** semantics. It's
the reference for anything beyond a single authenticator.

{{% notice tip %}}
Source: [`composed-auth/`](https://github.com/go-swagger/examples/tree/master/composed-auth).
Generated with `swagger generate server --name multi-auth-example --spec ./swagger.yml --principal models.Principal`.
The `restapi/configure_*.go` and `auth/authorizers.go` files are hand-written.
{{% /notice %}}

## Four schemes

The spec declares basic auth (`isRegistered`), an API key in either the header
(`isReseller`) or a query param (`isResellerQuery`), and a scoped `oauth2`-typed
scheme (`hasRole`) used purely to carry JWT scopes — go-swagger does not run an
OAuth2 flow here, it just extracts the required scopes and hands them to your
authorizer:

{{< code file="composed-auth/swagger.yml" lang="yaml" region="schemes" >}}

## Composing requirements with AND / OR

A `security` block on an operation is a **list of alternatives (OR)**, and each
alternative is a **map of schemes that must all pass (AND)**. So `/order/add`
accepts a registered customer, *or* a reseller (by header), *or* a reseller (by
query) — each combined with the right JWT role:

{{< code file="composed-auth/swagger.yml" lang="yaml" region="composed-security" >}}

An empty `security: []` on an operation opts it out entirely (public), overriding
the spec's top-level default.

## The authorizers

Each scheme maps to a hook whose signature depends on its type: basic auth gets
`(user, pass)`, an API key gets `(token)`, and a scoped scheme gets
`(token, scopes)`. The example delegates each to a function in
`auth/authorizers.go`:

{{< code file="composed-auth/restapi/configure_multi_auth_example.go" lang="go" region="wiring" >}}

The basic-auth authorizer is a plain credential check returning a typed
principal:

{{< code file="composed-auth/auth/authorizers.go" lang="go" region="is-registered" >}}

The scoped `HasRole` authorizer goes further: it parses the JWT, then checks the
token's claimed roles against the `scopes` the runtime passed in — the mechanism
that makes `hasRole: [ customer ]` in the spec actually mean something.

## Trying it

Generate test keys and JWTs, then exercise the composed requirements:

```shellsession
$ cd hack/tools && go run . gen-tokens        # RSA keypair + role JWTs
$ ./composed-auth/exerciser.sh                # sends a sequence of curl requests
```

## Related

- [Basic & API-key auth](../basic-and-apikey/) — the single-scheme starting point.
- [OAuth2 access-code](../oauth2-access-code/) — a real OAuth2 handshake (vs. JWT-scope extraction here).
