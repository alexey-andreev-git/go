package controller

import (
	"net/http"

	"what-to.com/internal/service"
)

func EntityHandler(w http.ResponseWriter, r *http.Request) {
	result := service.EntityServiceFunction(r)
	w.Write([]byte(result))
}
