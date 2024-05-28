package controller

import (
	"fmt"
	"io/fs"
	"net/http"

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
		httpHandlers: make(HttpHandlersT),
		appRes:       resources.NewAppSources(),
		config:       appConfig,
	}
	c.httpHandlers["/"] = ControllerHandlerT{
		Method:  "GET",
		Handler: c.RootHandler,
	}
	return c
}

func (c *FrontController) RootHandler(w http.ResponseWriter, r *http.Request) {
	c.config.GetLogger().Info(fmt.Sprintf("Start HTTP handler called from: %s, method: %s, path: %s", r.RemoteAddr, r.Method, r.URL.Path))
	subFS, err := fs.Sub(c.appRes.GetRes(), "appfs/frontend")
	if err != nil {
		c.config.GetLogger().Error("Error reading embedded FS for HTTP", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	fileServer := http.FileServer(http.FS(subFS))
	http.StripPrefix("/", fileServer).ServeHTTP(w, r)
	c.config.GetLogger().Info(fmt.Sprintf("[HTTP] [%s] [%s] served frontend files.", r.Method, r.URL.Path))
	c.config.GetLogger().Info(fmt.Sprintf("Finish HTTP handler called from: %s, method: %s, path: %s", r.RemoteAddr, r.Method, r.URL.Path))
}

func (c *FrontController) GetHandlers() HttpHandlersT {
	return c.httpHandlers
}
