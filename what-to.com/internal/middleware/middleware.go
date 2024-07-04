package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Middleware interface {
	GetHandler() http.Handler
	GetMiddlewareRouter() *mux.Router
}
