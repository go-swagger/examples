---
title: "Custom middleware"
weight: 4
description: "Plug classic net/http middleware into a server"
---

A generated server is a standard `net/http` stack, so any `http.Handler`
middleware composes with it. The **middleware** example wires two real concerns —
security response headers and Prometheus metrics — into a generated server using
only the hook points codegen already leaves for you. **No `--exclude-main` and no
custom router required.**

{{% notice tip %}}
Source: [`middleware/`](https://github.com/go-swagger/examples/tree/master/middleware).
Generated with `swagger generate server -A Greeter -f ./swagger.yml`. Answers
go-swagger issues [#2683](https://github.com/go-swagger/go-swagger/issues/2683)
(security headers) and [#1120](https://github.com/go-swagger/go-swagger/issues/1120)
(a `/metrics` endpoint).
{{% /notice %}}

## Two extension points

Codegen leaves two hooks in `restapi/configure_*.go`, and they run at different
stages of the request:

| Hook | Runs | Sees the matched route? | Use it for |
|------|------|:-----------------------:|-----------|
| `setupGlobalMiddleware` | Before swagger routing — wraps *everything* (spec, UI, all routes) | No | Cross-cutting concerns: security headers, panic recovery, a metrics mount |
| `setupMiddlewares` | After routing — only matched operations | Yes (`middleware.MatchedRouteFrom`) | Per-route concerns: instrumentation labelled by route template |

## The global hook — headers + metrics mount

`setupGlobalMiddleware` wraps the entire server. Reading outermost to innermost:
`metrics.Mount` intercepts `GET /metrics` so scrape traffic bypasses routing;
`unrolled/secure` adds HSTS and other headers to every response, including the
spec and UI:

{{< code file="middleware/restapi/configure_greeter.go" lang="go" region="setup-global" >}}

## The per-route hook — instrumentation

`setupMiddlewares` runs *after* routing, so the matched route is available. That's
what lets metrics be labelled by route:

{{< code file="middleware/restapi/configure_greeter.go" lang="go" region="setup-middlewares" >}}

## The go-swagger-specific glue: the route label

The one piece that's specific to a go-swagger server is the metrics `route`
label. It uses `middleware.MatchedRouteFrom(r).PathPattern` — the swagger path
*template* (`/greet/{name}`) rather than the literal path (`/greet/alice`) — so
Prometheus label cardinality stays bounded instead of exploding one series per
distinct URL:

{{< code file="middleware/internal/metrics/metrics.go" lang="go" region="instrument" >}}

## Trying it

```shellsession
$ go run ./middleware/cmd/greeter-server --port 8080
$ curl -i http://127.0.0.1:8080/api/greet
HTTP/1.1 200 OK
Strict-Transport-Security: max-age=63072000; includeSubDomains
X-Frame-Options: DENY
X-Content-Type-Options: nosniff
...
{"message":"hello"}

$ curl http://127.0.0.1:8080/metrics
http_requests_total{code="200",method="GET",route="/api/greet/{name}"} 3
```

The `route` label carries the `basePath` (`/api`) because the matched template is
the full path the router serves.

## Other middleware, same pattern

Any `http.Handler` middleware composes identically via `setupGlobalMiddleware` —
panic recovery (`gorilla/handlers.RecoveryHandler`), a request ID, an access log,
and so on. The example keeps to headers + metrics to stay focused on the wiring.

## Related

- [Todo list server](../../servers/todo-list/) — where the `configure_*.go` hooks come from.
- [Generation flags](../generation-flags/) — `--exclude-spec` and flag strategies (no `--exclude-main` needed here).
