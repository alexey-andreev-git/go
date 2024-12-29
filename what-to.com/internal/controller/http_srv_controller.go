package controller

import (
	"io/fs"
	"net/http"
	"os"
	"strings"

	"what-to.com/internal/config"
	"what-to.com/internal/models"
	"what-to.com/internal/resources"
	"what-to.com/internal/service"
)

type HttpFrontendControllerV1 struct {
	httpHandlers HttpHandlersT
	config       *config.Config
	newService   func() service.Service
}

const (
	frontRoutesPath = "/front/routes"
	frontFormsPath  = "/front/forms"
)

func NewFrontendControllerV1(appConfig *config.Config, newSvc func() service.Service) *HttpFrontendControllerV1 {
	c := &HttpFrontendControllerV1{
		config:     appConfig,
		newService: newSvc,
	}
	// c.AddHandler(ControllerHandlerT{Method: "GET", Handler: c.FrontendGet, Path: restWildcardPath})
	c.httpHandlers = HttpHandlersT{
		{Method: "GET", Handler: c.FrontendRoutesGet, Path: apiV1Path + frontRoutesPath + restWildcardPath},
		{Method: "GET", Handler: c.FrontendFormsGet, Path: apiV1Path + frontFormsPath + restWildcardPath},
		{Method: "GET", Handler: c.FrontendGet, Path: restWildcardPath},
	}
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

func (c *HttpFrontendControllerV1) FrontendRoutesGet(w http.ResponseWriter, r *http.Request) {
	routesData := &models.Routes{}
	svc := c.newService().(*service.FrontService)
	f := svc.GetServiceFuncs()[service.FrontRoutesGet]
	if _, err := f.Handler(routesData); err != nil {
		ErrorHandler(c.config.GetLogger(), w, errorMessage, err, http.StatusBadRequest)
		return
	}
	SetResponseJson(w, routesData)
}

func (c *HttpFrontendControllerV1) FrontendFormsGet(w http.ResponseWriter, r *http.Request) {
	formData := &models.Controls{}
	svc := c.newService().(*service.FrontService)
	f := svc.GetServiceFuncs()[service.FrontFormDtaGet]
	if _, err := f.Handler(formData); err != nil {
		ErrorHandler(c.config.GetLogger(), w, errorMessage, err, http.StatusBadRequest)
		return
	}
	SetResponseJson(w, formData)
}

func (c *HttpFrontendControllerV1) GetHandlers() HttpHandlersT {
	return c.httpHandlers
}
