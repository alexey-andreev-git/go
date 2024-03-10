package controller

import (
	"net/http"

	"what-to.com/internal/service"
)

func CompaniesHandler(w http.ResponseWriter, r *http.Request) {
	result := service.CompaniesServiceFunction()
	w.Write([]byte(result))
}
