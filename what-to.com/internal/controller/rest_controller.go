package controller

import (
	"net/http"

	"what-to.com/internal/resources"
	"what-to.com/internal/service"
)

type RestController struct {
	httpHandlers HttpHandlersT
	appRes       *resources.AppSources
}

func NewRestController() *RestController {
	c := &RestController{
		httpHandlers: make(HttpHandlersT),
	}
	c.appRes = resources.NewAppSources()
	c.httpHandlers["/entity/{rest:.*}"] = ControllerHandlerT{
		Method:  "GET",
		Handler: c.entityHandler,
	}
	c.httpHandlers["/api/{rest:.*}"] = ControllerHandlerT{
		Method:  "GET",
		Handler: c.apiHandler,
	}
	return c
}

func (c *RestController) entityHandler(w http.ResponseWriter, r *http.Request) {
	result := service.EntityServiceFunction(r)
	w.Write([]byte(result))
}

func (c *RestController) apiHandler(w http.ResponseWriter, r *http.Request) {
	result := service.EntityServiceFunction(r)
	w.Write([]byte(result))
}

func (c *RestController) GetHandlers() HttpHandlersT {
	return c.httpHandlers
}
