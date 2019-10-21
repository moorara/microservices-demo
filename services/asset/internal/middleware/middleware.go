package middleware

import "net/http"

// Middleware is a http middleware
type Middleware interface {
	Wrap(http.HandlerFunc) http.HandlerFunc
}
