---
title: "Basic & API-key auth"
weight: 1
description: "Basic auth and API-key security schemes"
---

This is the starting point for security in a generated server. A
`securityDefinition` in the spec becomes a generated **authenticator hook** you
implement, and — when you generate with a typed principal — every protected
handler receives that principal as a typed argument.

{{% notice tip %}}
Source: [`authentication/`](https://github.com/go-swagger/examples/tree/master/authentication).
Generated with `swagger generate server --name AuthSample --spec ./swagger.yml --principal models.Principal`.
{{% /notice %}}

## Declaring the scheme

The `authentication` example uses a single API-key scheme, carried in the
`x-token` header, and applies it to every endpoint via a top-level `security`
requirement:

{{< code file="authentication/swagger.yml" lang="yaml" region="security" >}}

## The generated authenticator hook

Because the spec names one `apiKey` scheme called `key`, the generated API exposes
an `api.KeyAuth` hook. The `--principal models.Principal` flag makes it return a
typed `*models.Principal`; the scaffold leaves it returning `NotImplemented`:

{{< code file="authentication/restapi/configure_auth_sample.go" lang="go" region="keyauth" >}}

You replace the body with your token check. On success return a principal; on
failure return an `errors.New(401, …)`:

```go
api.KeyAuth = func(token string) (*models.Principal, error) {
    if token == "abcdefuvwxyz" {
        prin := models.Principal(token)
        return &prin, nil
    }
    return nil, errors.New(401, "incorrect api key auth")
}
```

The returned principal is then passed to every handler protected by this scheme:

```go
api.CustomersGetIDHandler = customers.GetIDHandlerFunc(
    func(params customers.GetIDParams, principal *models.Principal) middleware.Responder {
        // principal is the value your KeyAuth hook returned
        ...
    })
```

## Basic auth is the same shape

A `type: basic` scheme works identically, except the runtime decodes the
`Authorization: Basic` header for you and the hook receives a **username/password
pair** instead of a single token:

```go
api.MyBasicAuth = func(user, pass string) (*models.Principal, error) { ... }
```

For a worked basic-auth authenticator — plus mixing several schemes — see
[Composed auth](../composed-auth/).

## Trying it

```shellsession
$ curl -i -H 'X-Token: abcdefuvwxyz' http://127.0.0.1:35307/api/customers
HTTP/1.1 501 Not Implemented          # authenticated, handler not implemented

$ curl -i -H 'X-Token: wrong' http://127.0.0.1:35307/api/customers
HTTP/1.1 401 Unauthorized
{"code":401,"message":"incorrect api key auth"}
```

## Related

- [Composed auth](../composed-auth/) — basic + API-key + scoped tokens, composed with AND/OR.
- [OAuth2 access-code](../oauth2-access-code/) — a full OAuth2 handshake.
- Hand-wiring auth without codegen? See the
  [runtime auth examples](https://go-openapi.github.io/runtime/).
