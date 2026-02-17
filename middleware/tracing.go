package middleware

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// TracingMiddleware wraps an HTTP handler with OpenTelemetry tracing
func TracingMiddleware(handler http.Handler, operationName string) http.Handler {
	return otelhttp.NewHandler(handler, operationName,
		otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
			return r.Method + " " + r.URL.Path
		}),
	)
}

// WrapHandler is a convenience function to wrap handlers with tracing
func WrapHandler(pattern string, handler http.Handler) (string, http.Handler) {
	return pattern, TracingMiddleware(handler, pattern)
}
