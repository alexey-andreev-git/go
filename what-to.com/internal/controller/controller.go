package controller

import (
	"net/http"
)

type (
	// ControllerHandlerT is a struct with Method and Handler for http requests
	ControllerHandlerT struct {
		Method  string
		Handler http.HandlerFunc
		// config  *config.Config
	}
	// HttpHandlersT is a map of path string to ControllerHandlerT
	HttpHandlersT map[string]ControllerHandlerT
	Controller    interface {
		GetHandlers() HttpHandlersT
	}
)
