package controller

import (
	"io/fs"
	"log"
	"net/http"

	"what-to.com/internal/config"
	"what-to.com/internal/resources"
)

type FrontController struct {
	httpHandlers HttpHandlersT
	appRes       *resources.AppSources
	config       *config.Config
}

func NewFrontController() *FrontController {
	c := &FrontController{
		httpHandlers: make(HttpHandlersT),
	}
	c.appRes = resources.NewAppSources()
	c.httpHandlers["/"] = ControllerHandlerT{
		Method:  "GET",
		Handler: c.RootHandler,
	}
	return c
}

func (c *FrontController) RootHandler(w http.ResponseWriter, r *http.Request) {
	c.appRes = resources.NewAppSources()
	// Create a sub-filesystem that points to 'appfs/frontend'.
	subFS, err := fs.Sub(c.appRes.GetRes(), "appfs/frontend")
	if err != nil {
		log.Fatal("Error readin embeded FS for HTTP", err)
	}
	fs := http.FS(subFS) // Convert embed.FS to http.FS
	fileServer := http.FileServer(fs)
	http.StripPrefix("/", fileServer).ServeHTTP(w, r)
	c.config.GetLogger().Info("[HTTP] [%s] [%s] served frontend files. Method: " + r.Method + " Path: " + r.URL.Path)
}

func (c *FrontController) GetHandlers() HttpHandlersT {
	return c.httpHandlers
}
