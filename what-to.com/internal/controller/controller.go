package controller

import (
	"net/http"
)

type (
	// ControllerHandlerT is a struct with Method and Handler for http requests
	ControllerHandlerT struct {
		Method  string
		Handler http.HandlerFunc
		Path    string
	}
	// HttpHandlersT is a slice of ControllerHandlerT
	HttpHandlersT []ControllerHandlerT
	Controller    interface {
		GetHandlers() HttpHandlersT
	}
)
