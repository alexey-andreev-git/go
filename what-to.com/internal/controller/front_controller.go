package controller

import (
	"io/fs"
	"net/http"
	"os"
	"strings"

	"what-to.com/internal/config"
	"what-to.com/internal/resources"
)

type FrontController struct {
	httpHandlers HttpHandlersT
	appRes       *resources.AppSources
	config       *config.Config
}

func NewFrontController(appConfig *config.Config) *FrontController {
	c := &FrontController{
		appRes: resources.NewAppSources(),
		config: appConfig,
	}
	c.httpHandlers = HttpHandlersT{
		ControllerHandlerT{Method: "GET", Handler: c.RootHandler, Path: "/{rest:.*}"},
	}
	return c
}

func (c *FrontController) RootHandler(w http.ResponseWriter, r *http.Request) {
	subFS, err := fs.Sub(c.appRes.GetRes(), "appfs/frontend")
	if err != nil {
		c.config.GetLogger().Error("Error reading embedded FS for HTTP", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	fileServer := http.FileServer(http.FS(subFS))
	path := r.URL.Path
	if path == "/" {
		path = "/index.html"
	}
	fullPath := path
	fullPath = strings.TrimPrefix(path, "/")
	if _, err := fs.Stat(subFS, fullPath); err != nil {
		if os.IsNotExist(err) {
			// If path doesn't exist then returning index.html
			indexFile, err := subFS.Open("index.html")
			if err != nil {
				c.config.GetLogger().Error("Error opening index.html from embedded FS", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			defer indexFile.Close()

			// Reading index.html
			indexData, err := fs.ReadFile(subFS, "index.html")
			if err != nil {
				c.config.GetLogger().Error("Error reading index.html from embedded FS", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Set header and write content of index.html to response
			w.Header().Set("Content-Type", "text/html")
			w.Write(indexData)
			return
		} else {
			// If another error then returning server error
			c.config.GetLogger().Error("Error reading file from embedded FS", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
	http.StripPrefix("/", fileServer).ServeHTTP(w, r)
}

func (c *FrontController) GetHandlers() HttpHandlersT {
	return c.httpHandlers
}
