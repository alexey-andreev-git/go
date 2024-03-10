package router

import (
	"what-to.com/internal/controller"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/entity", controller.EntityHandler).Methods("GET")
	return r
}
