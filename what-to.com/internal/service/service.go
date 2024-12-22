package service

import (
	"encoding/json"
	"net/http"

	"what-to.com/internal/logger"
)

const (
	apiV1Path        = "/api/v1"
	restWildcardPath = "/{rest:.*}"
	errorMessage     = "Error processing request"
)

type (
	RequestType int

	ServiceFunc struct {
		Handler func(map[string]interface{}) ([]byte, error)
		Method  string
		Path    string
	}

	Service interface {
		GetServiceFuncs() map[RequestType]ServiceFunc
		ServiceFunction(http.ResponseWriter, *http.Request, string, RequestType)
	}
)

func GetRequestBodyJson(w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
	defer r.Body.Close()
	bodyJson := make(map[string]interface{})
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&bodyJson)
	if err != nil {
		return nil, err
	}
	return bodyJson, nil
}

func ErrorHandler(l logger.Logger, w http.ResponseWriter, message string, err error, status int) {
	l.Error(message, err)
	http.Error(w, message+": "+err.Error(), status)
}
