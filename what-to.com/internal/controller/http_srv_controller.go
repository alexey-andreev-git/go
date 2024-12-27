package controller

import (
	"io/fs"
	"net/http"
	"os"
	"strings"

	"what-to.com/internal/config"
	"what-to.com/internal/repository"
	"what-to.com/internal/resources"
)

type HttpFrontendControllerV1 struct {
	httpHandlers HttpHandlersT
	config       *config.Config
}

func NewFrontendControllerV1(appConfig *config.Config, repo repository.Repository) *HttpFrontendControllerV1 {
	c := &HttpFrontendControllerV1{
		config: appConfig,
	}
	c.AddHandler(ControllerHandlerT{Method: "GET", Handler: c.FrontendGet, Path: restWildcardPath})
	return c
}

func (c *HttpFrontendControllerV1) FrontendGet(w http.ResponseWriter, r *http.Request) {
	res := resources.NewAppSources()
	subFS, err := fs.Sub(res.GetFs(), "appfs/frontend")
	if err != nil {
		ErrorHandler(c.config.GetLogger(), w, errorMessage, err, http.StatusInternalServerError)
		return
	}
	fileServer := http.FileServer(http.FS(subFS))
	if _, err := fs.Stat(subFS, strings.TrimPrefix(r.URL.Path, "/")); err != nil {
		if os.IsNotExist(err) {
			r.URL.Path = "/index.html"
		}
	}
	http.StripPrefix("/", fileServer).ServeHTTP(w, r)
}

func (c *HttpFrontendControllerV1) AddHandler(handler ControllerHandlerT) {
	c.httpHandlers = append(c.httpHandlers, ControllerHandlerT{Method: handler.Method, Handler: handler.Handler, Path: handler.Path})
}

func (c *HttpFrontendControllerV1) AddHandlers(handlers ...ControllerHandlerT) {
	for _, handler := range handlers {
		c.AddHandler(handler)
	}
}

func (c *HttpFrontendControllerV1) GetHandlers() HttpHandlersT {
	return c.httpHandlers
}
