package service

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"what-to.com/internal/resources"
)

// Import your_project_name/internal/repository here

func EntityServiceFunction(r *http.Request) string {
	// Here you would call your repository functions and implement business logic

	muxVars := mux.Vars(r)
	rest := muxVars["rest"]

	appRes := resources.NewAppSources()
	data, err := appRes.GetRes().ReadFile("appfs/sql/initdb.sql") // this is the embed.FS
	if err != nil {
		log.Fatalf("Ошибка при чтении файла: %v", err)
	}

	// Example: return r *http.Request as a string
	return ("Result: the entity\n" + r.RequestURI + "\n" + rest + "\n" + string(data))

	// return "Result: the entity"
}
