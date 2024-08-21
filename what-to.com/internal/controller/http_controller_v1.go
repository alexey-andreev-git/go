package controller

import (
	"net/http"

	"what-to.com/internal/config"
	"what-to.com/internal/service"
)

type HttpControllerV1 struct {
	httpHandlers HttpHandlersT
	config       *config.Config
	appService   service.Service
}

func NewHttpControllerV1(appConfig *config.Config, appService service.Service) *HttpControllerV1 {
	c := &HttpControllerV1{
		config:     appConfig,
		appService: appService,
	}
	c.registerHandlers()
	return c
}

func (c *HttpControllerV1) registerHandlers() {
	for rt, sf := range c.appService.GetServiceFuncs() {
		c.httpHandlers = append(c.httpHandlers, ControllerHandlerT{Method: sf.Method, Handler: c.handleRequest(rt), Path: sf.Path})
	}
}

// handleRequest is a generic handler for different types of requests.
func (c *HttpControllerV1) handleRequest(requestType service.RequestType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.appService.ServiceFunction(w, r, "1", requestType)
	}
}

func (c *HttpControllerV1) GetHandlers() HttpHandlersT {
	return c.httpHandlers
}
