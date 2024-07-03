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

const (
	apiV1Path = "/api/v1"
)

func NewRestController(appConfig *config.Config, appService service.Service) *RestController {
	c := &RestController{
		httpHandlers: make(HttpHandlersT),
		config:       appConfig,
		appService:   appService,
	}
	c.httpHandlers = HttpHandlersT{
		"entity_get":                     ControllerHandlerT{Method: "GET", Handler: c.entityV1GetHandler, Path: apiV1Path + "/entity/{rest:.*}"},
		"entity_post":                    ControllerHandlerT{Method: "POST", Handler: c.entityV1PostHandler, Path: apiV1Path + "/entity/{rest:.*}"},
		"entity_put":                     ControllerHandlerT{Method: "PUT", Handler: c.entityV1PutHandler, Path: apiV1Path + "/entity/{rest:.*}"},
		"entity_delete":                  ControllerHandlerT{Method: "DELETE", Handler: c.entityV1DeleteHandler, Path: apiV1Path + "/entity/{rest:.*}"},
		"entities_data_get":              ControllerHandlerT{Method: "GET", Handler: c.entityDataV1GetHandler, Path: apiV1Path + "/entities_data/{rest:.*}"},
		"entities_data_post":             ControllerHandlerT{Method: "POST", Handler: c.entityDataV1PostHandler, Path: apiV1Path + "/entities_data/{rest:.*}"},
		"entities_data_put":              ControllerHandlerT{Method: "PUT", Handler: c.entityDataV1PutHandler, Path: apiV1Path + "/entities_data/{rest:.*}"},
		"entities_data_delete":           ControllerHandlerT{Method: "DELETE", Handler: c.entityDataV1DeleteHandler, Path: apiV1Path + "/entities_data/{rest:.*}"},
		"entities_data_reference_get":    ControllerHandlerT{Method: "GET", Handler: c.entityDataRefV1GetHandler, Path: apiV1Path + "/entities_data_reference/{rest:.*}"},
		"entities_data_reference_post":   ControllerHandlerT{Method: "POST", Handler: c.entityDataRefV1PostHandler, Path: apiV1Path + "/entities_data_reference/{rest:.*}"},
		"entities_data_reference_put":    ControllerHandlerT{Method: "PUT", Handler: c.entityDataRefV1PutHandler, Path: apiV1Path + "/entities_data_reference/{rest:.*}"},
		"entities_data_reference_delete": ControllerHandlerT{Method: "DELETE", Handler: c.entityDataRefV1DeleteHandler, Path: apiV1Path + "/entities_data_reference/{rest:.*}"},
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

func (c *RestController) entityDataV1GetHandler(w http.ResponseWriter, r *http.Request) {
	c.config.GetLogger().Info(fmt.Sprintf("Start entity data handler called from: %s, method: %s, path: %s", r.RemoteAddr, r.Method, r.URL.Path))
	c.appService.EntityDataServiceFunctionGet(w, r, "1")
	c.config.GetLogger().Info(fmt.Sprintf("Finish entity data handler called from: %s, method: %s, path: %s", r.RemoteAddr, r.Method, r.URL.Path))
}

func (c *RestController) entityDataV1PostHandler(w http.ResponseWriter, r *http.Request) {
	c.config.GetLogger().Info(fmt.Sprintf("Start entity data handler called from: %s, method: %s, path: %s", r.RemoteAddr, r.Method, r.URL.Path))
	c.appService.EntityDataServiceFunctionPost(w, r, "1")
	c.config.GetLogger().Info(fmt.Sprintf("Finish entity data handler called from: %s, method: %s, path: %s", r.RemoteAddr, r.Method, r.URL.Path))
}

func (c *RestController) entityDataV1PutHandler(w http.ResponseWriter, r *http.Request) {
	c.config.GetLogger().Info(fmt.Sprintf("Start entity data handler called from: %s, method: %s, path: %s", r.RemoteAddr, r.Method, r.URL.Path))
	c.appService.EntityDataServiceFunctionPut(w, r, "1")
	c.config.GetLogger().Info(fmt.Sprintf("Finish entity data handler called from: %s, method: %s, path: %s", r.RemoteAddr, r.Method, r.URL.Path))
}

func (c *RestController) entityDataV1DeleteHandler(w http.ResponseWriter, r *http.Request) {
	c.config.GetLogger().Info(fmt.Sprintf("Start entity data handler called from: %s, method: %s, path: %s", r.RemoteAddr, r.Method, r.URL.Path))
	c.appService.EntityDataServiceFunctionDelete(w, r, "1")
	c.config.GetLogger().Info(fmt.Sprintf("Finish entity data handler called from: %s, method: %s, path: %s", r.RemoteAddr, r.Method, r.URL.Path))
}

func (c *RestController) entityDataRefV1GetHandler(w http.ResponseWriter, r *http.Request) {
	c.config.GetLogger().Info(fmt.Sprintf("Start entity data reference handler called from: %s, method: %s, path: %s", r.RemoteAddr, r.Method, r.URL.Path))
	c.appService.EntityDataRefServiceFunctionGet(w, r, "1")
	c.config.GetLogger().Info(fmt.Sprintf("Finish entity data reference handler called from: %s, method: %s, path: %s", r.RemoteAddr, r.Method, r.URL.Path))
}

func (c *RestController) entityDataRefV1PostHandler(w http.ResponseWriter, r *http.Request) {
	c.config.GetLogger().Info(fmt.Sprintf("Start entity data reference handler called from: %s, method: %s, path: %s", r.RemoteAddr, r.Method, r.URL.Path))
	c.appService.EntityDataRefServiceFunctionPost(w, r, "1")
	c.config.GetLogger().Info(fmt.Sprintf("Finish entity data reference handler called from: %s, method: %s, path: %s", r.RemoteAddr, r.Method, r.URL.Path))
}

func (c *RestController) entityDataRefV1PutHandler(w http.ResponseWriter, r *http.Request) {
	c.config.GetLogger().Info(fmt.Sprintf("Start entity data reference handler called from: %s, method: %s, path: %s", r.RemoteAddr, r.Method, r.URL.Path))
	c.appService.EntityDataRefServiceFunctionPut(w, r, "1")
	c.config.GetLogger().Info(fmt.Sprintf("Finish entity data reference handler called from: %s, method: %s, path: %s", r.RemoteAddr, r.Method, r.URL.Path))
}

func (c *RestController) entityDataRefV1DeleteHandler(w http.ResponseWriter, r *http.Request) {
	c.config.GetLogger().Info(fmt.Sprintf("Start entity data reference handler called from: %s, method: %s, path: %s", r.RemoteAddr, r.Method, r.URL.Path))
	c.appService.EntityDataRefServiceFunctionDelete(w, r, "1")
	c.config.GetLogger().Info(fmt.Sprintf("Finish entity data reference handler called from: %s, method: %s, path: %s", r.RemoteAddr, r.Method, r.URL.Path))
}

func (c *RestController) GetHandlers() HttpHandlersT {
	return c.httpHandlers
}
