package controller

import (
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
	apiV1Path         = "/api/v1"
	restWildcardPath  = "/{rest:.*}"
	entityPath        = "/entity"
	entityDataPath    = "/entities_data"
	entityDataRefPath = "/entities_data_reference"
)

func NewRestController(appConfig *config.Config, appService service.Service) *RestController {
	c := &RestController{
		config:     appConfig,
		appService: appService,
	}
	c.registerHandlers()
	return c
}

func (c *RestController) registerHandlers() {
	c.httpHandlers = HttpHandlersT{
		{Method: "GET", Handler: c.handleRequest(service.EntityGet), Path: apiV1Path + entityPath + restWildcardPath},
		{Method: "POST", Handler: c.handleRequest(service.EntityPost), Path: apiV1Path + entityPath + restWildcardPath},
		{Method: "PUT", Handler: c.handleRequest(service.EntityPut), Path: apiV1Path + entityPath + restWildcardPath},
		{Method: "DELETE", Handler: c.handleRequest(service.EntityDelete), Path: apiV1Path + entityPath + restWildcardPath},
		{Method: "GET", Handler: c.handleRequest(service.EntityDataGet), Path: apiV1Path + entityDataPath + restWildcardPath},
		{Method: "POST", Handler: c.handleRequest(service.EntityDataPost), Path: apiV1Path + entityDataPath + restWildcardPath},
		{Method: "PUT", Handler: c.handleRequest(service.EntityDataPut), Path: apiV1Path + entityDataPath + restWildcardPath},
		{Method: "DELETE", Handler: c.handleRequest(service.EntityDataDelete), Path: apiV1Path + entityDataPath + restWildcardPath},
		{Method: "GET", Handler: c.handleRequest(service.EntityDataRefGet), Path: apiV1Path + entityDataRefPath + restWildcardPath},
		{Method: "POST", Handler: c.handleRequest(service.EntityDataRefPost), Path: apiV1Path + entityDataRefPath + restWildcardPath},
		{Method: "PUT", Handler: c.handleRequest(service.EntityDataRefPut), Path: apiV1Path + entityDataRefPath + restWildcardPath},
		{Method: "DELETE", Handler: c.handleRequest(service.EntityDataRefDelete), Path: apiV1Path + entityDataRefPath + restWildcardPath},
	}
}

// handleRequest is a generic handler for different types of requests.
func (c *RestController) handleRequest(requestType service.RequestType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.appService.EntityServiceFunction(w, r, "1", requestType)
	}
}

func (c *RestController) GetHandlers() HttpHandlersT {
	return c.httpHandlers
}
