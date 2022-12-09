package middleware

import (
	"net/http"
)

type (
	middlewareOptions func(http.Handler) http.Handler
	Middleware        []middlewareOptions
)

func New(opt ...middlewareOptions) Middleware {
	var m Middleware

	return append(m, opt...)
}

func (m Middleware) Configure(handler http.Handler) http.Handler {
	if handler == nil {
		handler = http.DefaultServeMux
	}

	for i := range m {
		handler = m[len(m)-1-i](handler)
	}

	return handler
}
