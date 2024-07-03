package controller

import (
	"fmt"
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
		httpHandlers: make(HttpHandlersT),
		appRes:       resources.NewAppSources(),
		config:       appConfig,
	}
	c.httpHandlers = HttpHandlersT{
		"front_root": ControllerHandlerT{Method: "GET", Handler: c.RootHandler, Path: "/{rest:.*}"},
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
	path := r.URL.Path
	if path == "/" {
		path = "/index.html"
	}
	fullPath := path
	fullPath = strings.TrimPrefix(path, "/")
	if _, err := fs.Stat(subFS, fullPath); err != nil {
		if os.IsNotExist(err) {
			// Если файл не существует, возвращаем index.html
			// r.URL.Path = "/index.html"
			indexFile, err := subFS.Open("index.html")
			if err != nil {
				c.config.GetLogger().Error("Error opening index.html from embedded FS", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			defer indexFile.Close()

			// Читаем содержимое index.html
			indexData, err := fs.ReadFile(subFS, "index.html")
			if err != nil {
				c.config.GetLogger().Error("Error reading index.html from embedded FS", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Устанавливаем заголовки и пишем содержимое index.html в ответ
			w.Header().Set("Content-Type", "text/html")
			w.Write(indexData)
			return
		} else {
			// Если произошла другая ошибка, возвращаем ошибку сервера
			c.config.GetLogger().Error("Error reading file from embedded FS", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
	http.StripPrefix("/", fileServer).ServeHTTP(w, r)
	c.config.GetLogger().Info(fmt.Sprintf("[HTTP] [%s] [%s] served frontend files.", r.Method, r.URL.Path))
	c.config.GetLogger().Info(fmt.Sprintf("Finish HTTP handler called from: %s, method: %s, path: %s", r.RemoteAddr, r.Method, r.URL.Path))
}

func (c *FrontController) GetHandlers() HttpHandlersT {
	return c.httpHandlers
}
