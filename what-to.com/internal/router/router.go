package router

import (
	"io/fs"
	"log"
	"net/http"

	"what-to.com/internal/controller"
	"what-to.com/internal/resources"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		appRes := resources.NewAppSources()
		// Create a sub-filesystem that points to 'appfs/frontend'.
		subFS, err := fs.Sub(appRes.GetRes(), "appfs/frontend")
		if err != nil {
			log.Fatal("Error readin embeded FS for HTTP", err)
		}
		fs := http.FS(subFS) // Convert embed.FS to http.FS
		fileServer := http.FileServer(fs)
		http.StripPrefix("/", fileServer).ServeHTTP(w, r)
	}).Methods("GET")

	r.HandleFunc("/entity/{rest:.*}", controller.EntityHandler).Methods("GET")
	return r
}
