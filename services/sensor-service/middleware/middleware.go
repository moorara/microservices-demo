package middleware

import "net/http"

type (
	// Middleware represents a http middleware
	Middleware interface {
		Wrap(http.HandlerFunc) http.HandlerFunc
	}
)

// WrapAll wraps a http handler with a set of middleware
func WrapAll(handler http.HandlerFunc, middleware ...Middleware) http.HandlerFunc {
	wrappedHandler := handler
	for i := len(middleware) - 1; i >= 0; i-- {
		wrappedHandler = middleware[i].Wrap(wrappedHandler)
	}

	return wrappedHandler
}
