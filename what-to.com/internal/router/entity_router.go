package router

import (
	"what-to.com/internal/controller"

	"github.com/gorilla/mux"
)

type EntityRouter struct {
	muxRouter *mux.Router
	//map of controllers [entityName] => controller
	entityControllers map[string]controller.Controller
}

func NewEntityRouter() *EntityRouter {
	r := &EntityRouter{}
	r.entityControllers = make(map[string]controller.Controller)
	r.muxRouter = mux.NewRouter()
	return r
}

func (r *EntityRouter) AddController(n string, c controller.Controller) {
	r.entityControllers[n] = c

	for path, handler := range c.GetHandlers() {
		r.muxRouter.HandleFunc(path, handler.Handler).Methods(handler.Method)
	}

}

func (r *EntityRouter) GetMuxRouter() *mux.Router {
	return r.muxRouter
}
