package controller

import (
	"net/http"
)

type (
	// ControllerHandlerT is a struct with Method and Handler for http requests
	ControllerHandlerT struct {
		Handler http.HandlerFunc
		Method  string
		Path    string
	}
	// HttpHandlersT is a slice of ControllerHandlerT
	HttpHandlersT []ControllerHandlerT
	Controller    interface {
		GetHandlers() HttpHandlersT
	}
)
