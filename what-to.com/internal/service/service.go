package service

import (
	"net/http"
)

type (
	Service interface {
		// CreateEntity(entity.Entity) (*entity.Entity, error)
		// GetEntityById(int) (*entity.Entity, error)
		// GetEntityByName(string) (*entity.Entity, error)
		EntityServiceFunctionGet(http.ResponseWriter, *http.Request, string)
		EntityServiceFunctionPost(http.ResponseWriter, *http.Request, string)
		EntityServiceFunctionPut(http.ResponseWriter, *http.Request, string)
		EntityServiceFunctionDelete(http.ResponseWriter, *http.Request, string)
	}
)
