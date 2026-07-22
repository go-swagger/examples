// SPDX-FileCopyrightText: Copyright 2015-2026 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

// Package metrics is the small, go-swagger-specific glue around
// prometheus/client_golang used by the middleware example.
//
// The example uses the default global Prometheus registry to keep the focus
// on middleware wiring rather than on Prometheus best practices. A real
// service may prefer a private [prometheus.Registry].
package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-openapi/runtime/middleware"
)

const (
	labelMethod = "method"
	labelRoute  = "route"
	labelCode   = "code"
)

var (
	labels = []string{labelMethod, labelRoute, labelCode}

	requestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests served, partitioned by method, swagger route template and status code.",
	}, labels)

	requestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Latency of HTTP requests, partitioned by method, swagger route template and status code.",
		Buckets: prometheus.DefBuckets,
	}, labels)
)

// Mount short-circuits GET /metrics to the Prometheus exposition handler and
// delegates every other request to next.
//
// Place Mount in setupGlobalMiddleware so that /metrics bypasses the swagger
// router and any other request-level middleware (security headers,
// instrumentation, ...) that should not apply to scrape traffic.
func Mount(next http.Handler) http.Handler {
	promHandler := promhttp.Handler()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/metrics" {
			promHandler.ServeHTTP(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// snippet:instrument

// Instrument records request count and latency for the wrapped handler.
//
// It is meant to be installed in the generated server's setupMiddlewares hook
// (i.e. after swagger routing) so that [middleware.MatchedRouteFrom] returns
// the matched route. The route label is the swagger path template (e.g.
// "/greet/{name}") rather than the literal request path, to keep the metric
// label cardinality bounded.
func Instrument(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)

		route := "unmatched"
		if mr := middleware.MatchedRouteFrom(r); mr != nil && mr.PathPattern != "" {
			route = mr.PathPattern
		}

		obs := prometheus.Labels{
			labelMethod: r.Method,
			labelRoute:  route,
			labelCode:   strconv.Itoa(rec.status),
		}
		requestsTotal.With(obs).Inc()
		requestDuration.With(obs).Observe(time.Since(start).Seconds())
	})
}

// endsnippet:instrument

type statusRecorder struct {
	http.ResponseWriter

	status      int
	wroteHeader bool
}

func (s *statusRecorder) WriteHeader(code int) {
	if !s.wroteHeader {
		s.status = code
		s.wroteHeader = true
	}
	s.ResponseWriter.WriteHeader(code)
}
