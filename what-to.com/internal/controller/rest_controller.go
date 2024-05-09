package controller

import (
	"fmt"
	"net/http"

	"what-to.com/internal/config"
	"what-to.com/internal/resources"
	"what-to.com/internal/service"
)

type RestController struct {
	httpHandlers HttpHandlersT
	appRes       *resources.AppSources
	config       *config.Config
}

func NewRestController(appConfig *config.Config) *RestController {
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
	c.config = appConfig
	return c
}

func (c *RestController) entityHandler(w http.ResponseWriter, r *http.Request) {
	c.config.GetLogger().Info(fmt.Sprintf("Start entity handler called from: %s, method: %s, path: %s", r.RemoteAddr, r.Method, r.URL.Path))
	result := service.EntityServiceFunction(r, c.config)
	w.Write([]byte(result))
	c.config.GetLogger().Info(fmt.Sprintf("Finish entity handler called from: %s, method: %s, path: %s", r.RemoteAddr, r.Method, r.URL.Path))
}

func (c *RestController) apiHandler(w http.ResponseWriter, r *http.Request) {
	c.config.GetLogger().Info(fmt.Sprintf("Start api handler called from: %s, method: %s, path: %s", r.RemoteAddr, r.Method, r.URL.Path))
	result := service.EntityServiceFunction(r, c.config)
	w.Write([]byte(result))
	c.config.GetLogger().Info(fmt.Sprintf("Finish api handler called from: %s, method: %s, path: %s", r.RemoteAddr, r.Method, r.URL.Path))
}

func (c *RestController) GetHandlers() HttpHandlersT {
	return c.httpHandlers
}
