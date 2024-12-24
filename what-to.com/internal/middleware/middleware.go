package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
)

type (
	Middleware interface {
		Use()
		ChainMiddleware() http.Handler
		GetMiddlewareRouter() *mux.Router
	}
)
