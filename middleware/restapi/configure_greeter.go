// SPDX-FileCopyrightText: Copyright 2015-2026 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/unrolled/secure"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/go-swagger/examples/middleware/internal/metrics"
	"github.com/go-swagger/examples/middleware/models"
	"github.com/go-swagger/examples/middleware/restapi/operations"
	"github.com/go-swagger/examples/middleware/restapi/operations/greeter"
)

//go:generate swagger generate server --target ../../middleware --name Greeter --spec ../swagger.yml --principal any

func configureFlags(api *operations.GreeterAPI) {
	// api.CommandLineOptionsGroups = []cmdutils.CommandLineOptionsGroup{ ... }
	_ = api
}

func configureAPI(api *operations.GreeterAPI) http.Handler {
	api.ServeError = errors.ServeError

	api.UseSwaggerUI()

	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()

	api.GreeterGreetHandler = greeter.GreetHandlerFunc(func(_ greeter.GreetParams) middleware.Responder {
		return greeter.NewGreetOK().WithPayload(&models.Greeting{Message: stringPtr("hello")})
	})
	api.GreeterGreetNameHandler = greeter.GreetNameHandlerFunc(func(params greeter.GreetNameParams) middleware.Responder {
		msg := "hello, " + params.Name
		return greeter.NewGreetNameOK().WithPayload(&models.Greeting{Message: &msg})
	})

	api.PreServerShutdown = func() {}
	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

func stringPtr(s string) *string { return &s }

// configureTLS is called before the HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	_ = tlsConfig
}

// configureServer is called as soon as the server is initialized but not run yet.
// scheme is set to "http", "https" or "unix".
func configureServer(server *http.Server, scheme, addr string) {
	_ = server
	_ = scheme
	_ = addr
}

// setupMiddlewares wraps the swagger handler after routing.
//
// At this point, [middleware.MatchedRouteFrom] returns the matched route, so
// this is the right place to plug per-route instrumentation. Requests that do
// not match an operation (404/405) never reach this layer; that is by design
// for an example whose metrics describe API operation behaviour.
func setupMiddlewares(handler http.Handler) http.Handler {
	return metrics.Instrument(handler)
}

// setupGlobalMiddleware wraps everything the server serves, including the
// swagger spec document and the embedded UI.
//
// Order matters. From outermost to innermost:
//
//  1. metrics.Mount intercepts GET /metrics so scrape traffic bypasses the
//     swagger router as well as the headers and instrumentation below.
//  2. unrolled/secure injects HSTS and other security response headers on
//     every response, including /swagger.json and the UI assets.
//  3. The swagger handler itself (which will pass through setupMiddlewares
//     once a route is matched).
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	sec := secure.New(secure.Options{
		STSSeconds:           63072000,
		STSIncludeSubdomains: true,
		// ForceSTSHeader: true emits HSTS even on plain-HTTP requests so the
		// example is easy to exercise with `curl http://...`. In production
		// you would either serve over HTTPS directly (the library detects it
		// from the TLS connection) or set SSLProxyHeaders so detection works
		// behind a TLS-terminating proxy — and leave ForceSTSHeader at false.
		ForceSTSHeader:     true,
		FrameDeny:          true,
		ContentTypeNosniff: true,
		ReferrerPolicy:     "no-referrer",
		// Send "X-XSS-Protection: 0" explicitly. The default BrowserXssFilter
		// option sends "1; mode=block", which is now considered harmful in
		// some browsers; modern OWASP / MDN guidance is to disable the
		// legacy XSS auditor outright. See README.
		CustomBrowserXssValue: "0",
	})

	return metrics.Mount(sec.Handler(handler))
}
