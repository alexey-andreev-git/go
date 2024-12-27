package controller

import (
	"net/http"

	"what-to.com/internal/logger"
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
		AddHandler(handler ControllerHandlerT)
		AddHandlers(handlers ...ControllerHandlerT)
		GetHandlers() HttpHandlersT
	}
)

const (
	restWildcardPath = "/{rest:.*}"
	errorMessage     = "Error processing request"
)

func ErrorHandler(l logger.Logger, w http.ResponseWriter, message string, err error, status int) {
	l.Error(message, err)
	http.Error(w, message+": "+err.Error(), status)
}
