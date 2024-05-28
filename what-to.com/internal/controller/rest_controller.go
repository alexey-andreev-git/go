package controller

import (
	"fmt"
	"net/http"

	"what-to.com/internal/config"
	"what-to.com/internal/service"
)

type RestController struct {
	httpHandlers HttpHandlersT
	config       *config.Config
	appService   service.Service
}

func NewRestController(appConfig *config.Config, appService service.Service) *RestController {
	c := &RestController{
		httpHandlers: make(HttpHandlersT),
		config:       appConfig,
		appService:   appService,
	}
	c.httpHandlers["entity_get"] = ControllerHandlerT{
		Method:  "GET",
		Handler: c.entityV1GetHandler,
		Path:    "/api/v1/entity/{rest:.*}",
	}
	c.httpHandlers["entity_post"] = ControllerHandlerT{
		Method:  "POST",
		Handler: c.entityV1PostHandler,
		Path:    "/api/v1/entity/{rest:.*}",
	}
	c.httpHandlers["entity_put"] = ControllerHandlerT{
		Method:  "PUT",
		Handler: c.entityV1PutHandler,
		Path:    "/api/v1/entity/{rest:.*}",
	}
	c.httpHandlers["entity_delete"] = ControllerHandlerT{
		Method:  "DELETE",
		Handler: c.entityV1DeleteHandler,
		Path:    "/api/v1/entity/{rest:.*}",
	}
	return c
}

func (c *RestController) entityV1GetHandler(w http.ResponseWriter, r *http.Request) {
	c.config.GetLogger().Info(fmt.Sprintf("Start entity handler called from: %s, method: %s, path: %s", r.RemoteAddr, r.Method, r.URL.Path))
	c.appService.EntityServiceFunctionGet(w, r, "1")
	c.config.GetLogger().Info(fmt.Sprintf("Finish entity handler called from: %s, method: %s, path: %s", r.RemoteAddr, r.Method, r.URL.Path))
}

func (c *RestController) entityV1PostHandler(w http.ResponseWriter, r *http.Request) {
	c.config.GetLogger().Info(fmt.Sprintf("Start API handler called from: %s, method: %s, path: %s", r.RemoteAddr, r.Method, r.URL.Path))
	c.appService.EntityServiceFunctionPost(w, r, "1")
	c.config.GetLogger().Info(fmt.Sprintf("Finish API handler called from: %s, method: %s, path: %s", r.RemoteAddr, r.Method, r.URL.Path))
}

func (c *RestController) entityV1PutHandler(w http.ResponseWriter, r *http.Request) {
	c.config.GetLogger().Info(fmt.Sprintf("Start API handler called from: %s, method: %s, path: %s", r.RemoteAddr, r.Method, r.URL.Path))
	c.appService.EntityServiceFunctionPut(w, r, "1")
	c.config.GetLogger().Info(fmt.Sprintf("Finish API handler called from: %s, method: %s, path: %s", r.RemoteAddr, r.Method, r.URL.Path))
}

func (c *RestController) entityV1DeleteHandler(w http.ResponseWriter, r *http.Request) {
	c.config.GetLogger().Info(fmt.Sprintf("Start API handler called from: %s, method: %s, path: %s", r.RemoteAddr, r.Method, r.URL.Path))
	c.appService.EntityServiceFunctionDelete(w, r, "1")
	c.config.GetLogger().Info(fmt.Sprintf("Finish API handler called from: %s, method: %s, path: %s", r.RemoteAddr, r.Method, r.URL.Path))
}

func (c *RestController) GetHandlers() HttpHandlersT {
	return c.httpHandlers
}
