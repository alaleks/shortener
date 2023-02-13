// Package middleware implements options for the web server.
package middleware

import (
	"net/http"
)

// Middleware instances.
type (
	middlewareOptions func(http.Handler) http.Handler
	Middleware        []middlewareOptions
)

// New initializes the middleware instances with options.
func New(opt ...middlewareOptions) Middleware {
	var m Middleware

	return append(m, opt...)
}

// Configure includes middleware`s configurations for routes.
func (m Middleware) Configure(handler http.Handler) http.Handler {
	if handler == nil {
		handler = http.DefaultServeMux
	}

	for i := range m {
		handler = m[len(m)-1-i](handler)
	}

	return handler
}
