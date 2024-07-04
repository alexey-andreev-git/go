package service

import "net/http"

type RequestType int

type (
	Service interface {
		EntityServiceFunction(http.ResponseWriter, *http.Request, string, RequestType)
	}
)

const (
	EntityGet RequestType = iota
	EntityPost
	EntityPut
	EntityDelete
	EntityDataGet
	EntityDataPost
	EntityDataPut
	EntityDataDelete
	EntityDataRefGet
	EntityDataRefPost
	EntityDataRefPut
	EntityDataRefDelete
)
