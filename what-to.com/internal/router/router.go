package router

import (
	"github.com/gorilla/mux"
	"what-to.com/internal/controller"
)

type Router interface {
	GetMuxRouter() *mux.Router
	AddController(string, controller.Controller)
}
