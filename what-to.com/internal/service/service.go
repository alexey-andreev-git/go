package service

import (
	"net/http"
)

type (
	Service interface {
		EntityServiceFunctionGet(http.ResponseWriter, *http.Request, string)
		EntityServiceFunctionPost(http.ResponseWriter, *http.Request, string)
		EntityServiceFunctionPut(http.ResponseWriter, *http.Request, string)
		EntityServiceFunctionDelete(http.ResponseWriter, *http.Request, string)

		EntityDataServiceFunctionGet(http.ResponseWriter, *http.Request, string)
		EntityDataServiceFunctionPost(http.ResponseWriter, *http.Request, string)
		EntityDataServiceFunctionPut(http.ResponseWriter, *http.Request, string)
		EntityDataServiceFunctionDelete(http.ResponseWriter, *http.Request, string)

		EntityDataRefServiceFunctionGet(http.ResponseWriter, *http.Request, string)
		EntityDataRefServiceFunctionPost(http.ResponseWriter, *http.Request, string)
		EntityDataRefServiceFunctionPut(http.ResponseWriter, *http.Request, string)
		EntityDataRefServiceFunctionDelete(http.ResponseWriter, *http.Request, string)
	}
)
