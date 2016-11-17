package router

import (
	"net/http"

	"github.com/gernest/utron/base"
)

//MiddlewareType is the kind of middleware. Utron support middleware with
//variary of signatures.
type MiddlewareType int

const (
	//PlainMiddleware is the middleware with signature
	// func(http.Handler)http.Handler
	PlainMiddleware MiddlewareType = iota

	//CtxMiddleware is the middlewate with signature
	// func(*base.Context)error
	CtxMiddleware
)

//Middleware is the utron middleware
type Middleware struct {
	Type  MiddlewareType
	value interface{}
}

//ToHandler returns a func(http.Handler) http.Handler from the Middleware. Utron
//uses alice to chain middleware.
//
// Use this method to get alice compatible middleware.
func (m *Middleware) ToHandler(ctx *base.Context) func(http.Handler) http.Handler {
	switch m.Type {
	case PlainMiddleware:
		return m.value.(func(http.Handler) http.Handler)
	case CtxMiddleware:
		fn := m.value.(func(*base.Context) error)
		return func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				err := fn(ctx)
				if err != nil {
					return
				}
				h.ServeHTTP(w, r)
			})
		}

	}
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		})
	}
}
