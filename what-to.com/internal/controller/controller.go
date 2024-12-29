package controller

import (
	"encoding/json"
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
		GetHandlers() HttpHandlersT
	}
)

const (
	apiV1Path        = "/api/v1"
	restWildcardPath = "/{rest:.*}"
	errorMessage     = "Error processing request"
)

func ErrorHandler(l logger.Logger, w http.ResponseWriter, message string, err error, status int) {
	l.Error(message, err)
	http.Error(w, message+": "+err.Error(), status)
}

func SetResponseJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}
}
