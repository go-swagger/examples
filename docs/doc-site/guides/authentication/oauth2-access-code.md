---
title: "OAuth2 access-code"
weight: 3
description: "A full OAuth2 access-code handshake"
---

The other auth examples *validate* a token the client already has. This one runs
the full **OAuth2 access-code handshake** — redirecting the user to an identity
provider (Google), receiving a callback, exchanging the code for a token, then
using that token to authenticate API calls.

{{% notice tip %}}
Source: [`oauth2/`](https://github.com/go-swagger/examples/tree/master/oauth2).
Generated with `swagger generate server --name oauthSample --spec ./swagger.yml --principal models.Principal`.
The handshake lives in the hand-written `restapi/implementation.go`.
{{% /notice %}}

## The scheme

The spec declares an `oauth2` scheme with the `accessCode` flow and its
authorization/token URLs. Unlike [composed-auth](../composed-auth/) — which only
borrows the `oauth2` type to carry scopes — here the URLs are real and drive an
actual handshake:

{{< code file="oauth2/swagger.yml" lang="yaml" region="security" >}}

{{% notice note %}}
go-swagger does not implement the OAuth2 workflow for you: the generator produces
the authenticator hook and the routing, but the redirect/callback/exchange dance
is application code. That's exactly what this example provides.
{{% /notice %}}

## Step 1 — redirect to the provider

The public `/login` endpoint sends the user to Google's consent screen, using the
`golang.org/x/oauth2` config built in `implementation.go`:

{{< code file="oauth2/restapi/implementation.go" lang="go" region="login" >}}

## Step 2 — handle the callback and exchange the code

Google redirects back to `/auth/callback` with a `state` and a `code`. The
handler verifies `state`, then exchanges the code for an access token via the
oauth2 client:

{{< code file="oauth2/restapi/implementation.go" lang="go" region="callback" >}}

## Step 3 — authenticate API calls with the token

Every protected endpoint runs through the generated `OauthSecurityAuth` hook. It
validates the bearer token (here by calling Google's userinfo endpoint) and
returns the principal — the token string itself in this minimal example:

{{< code file="oauth2/restapi/configure_oauth_sample.go" lang="go" region="oauth-auth" >}}

## Setup

Register an OAuth client at the [Google credentials console](https://console.cloud.google.com/apis/credentials/),
set the callback URL to `http://127.0.0.1:12345/api/auth/callback`, and put the
resulting client ID/secret into the `var` block of `implementation.go`. Then:

```shellsession
$ go run ./oauth2/cmd/oauth-sample-server/main.go --port 12345
# open http://127.0.0.1:12345/api/login in a browser, log in, copy the token
$ curl -i -H 'Authorization: Bearer <TOKEN>' http://127.0.0.1:12345/api/customers
```

## Related

- [Basic & API-key auth](../basic-and-apikey/) — validating a pre-shared credential.
- [Composed auth](../composed-auth/) — mixing schemes and extracting JWT scopes.
