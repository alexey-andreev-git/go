package service

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"what-to.com/internal/config"
	"what-to.com/internal/repository"
	"what-to.com/internal/resources"
)

// Struct of entity service
type FrontControllerService struct {
	appRepository repository.Repository
	appConfig     *config.Config
	appRes        *resources.AppSources
	serviceFuncs  map[RequestType]ServiceFunc
}

const (
	FrontGet RequestType = iota
	FrontPost
	FrontPut
	FrontDelete
)

func NewFrontService(appConfig *config.Config, appRepo repository.Repository) *FrontControllerService {
	s := &FrontControllerService{
		appConfig:     appConfig,
		appRepository: appRepo,
		appRes:        resources.NewAppSources(),
	}
	s.registerServiceFuncs()
	return s
}

func (s *FrontControllerService) registerServiceFuncs() {
	s.serviceFuncs = map[RequestType]ServiceFunc{
		EntityGet:    {s.V1FrontControllerServiceGet, "GET", restWildcardPath},
		EntityPost:   {s.V1FrontControllerServicePost, "POST", restWildcardPath},
		EntityPut:    {s.V1FrontControllerServicePut, "PUT", restWildcardPath},
		EntityDelete: {s.V1FrontControllerServiceDelete, "DELETE", restWildcardPath},
	}
}

// func (s *FrontControllerService) errorHandler(w http.ResponseWriter, message string, err error, status int) {
// 	s.appConfig.GetLogger().Error(message, err)
// 	http.Error(w, errorMessage+": "+err.Error(), status)
// }

func (s *FrontControllerService) GetRequestParamsJson(w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
	defer r.Body.Close()
	params := make(map[string]interface{})
	params["path"] = r.URL.Path
	return params, nil
}

func (s *FrontControllerService) GetServiceFuncs() map[RequestType]ServiceFunc {
	return s.serviceFuncs
}

func (s *FrontControllerService) ServiceFunction(w http.ResponseWriter, r *http.Request, ver string, reqType RequestType) {
	if ver != "1" {
		ErrorHandler(s.appConfig.GetLogger(), w, errorMessage, fmt.Errorf("incorrect request version"), http.StatusBadRequest)
		return
	}
	handlerFunc, ok := s.serviceFuncs[reqType]
	if !ok {
		ErrorHandler(s.appConfig.GetLogger(), w, errorMessage, fmt.Errorf("incorrect request type"), http.StatusBadRequest)
		return
	}
	jsonBody, err := s.GetRequestParamsJson(w, r)
	if err != nil {
		ErrorHandler(s.appConfig.GetLogger(), w, errorMessage, err, http.StatusInternalServerError)
		return
	}
	subFS, err := s.getFrontFs()
	if err != nil {
		ErrorHandler(s.appConfig.GetLogger(), w, errorMessage, err, http.StatusInternalServerError)
		return
	}
	fileServer := http.FileServer(http.FS(subFS))
	jsonBody["subFs"] = subFS
	respJson, rerr := handlerFunc.Handler(jsonBody)
	if rerr != nil {
		ErrorHandler(s.appConfig.GetLogger(), w, errorMessage, rerr, http.StatusBadRequest)
	}
	if respJson != nil {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write(respJson)
	} else {
		http.StripPrefix("/", fileServer).ServeHTTP(w, r)
	}
}

func (s *FrontControllerService) getFrontFs() (fs.FS, error) {
	subFS, err := fs.Sub(s.appRes.GetRes(), "appfs/frontend")
	if err != nil {
		return nil, err
	}
	return subFS, nil
}

func (s *FrontControllerService) V1FrontControllerServiceGet(jsonBody map[string]interface{}) ([]byte, error) {
	var response []byte = nil
	path := jsonBody["path"].(string)
	if path == "/" {
		path = "/index.html"
	}
	fullPath := path
	fullPath = strings.TrimPrefix(path, "/")
	subFS := jsonBody["subFs"].(fs.FS)
	if _, err := fs.Stat(subFS, fullPath); err != nil {
		if os.IsNotExist(err) {
			// If path doesn't exist then returning index.html
			indexFile, err := subFS.Open("index.html")
			if err != nil {
				return nil, err
			}
			defer indexFile.Close()

			// Reading index.html
			indexData, err := fs.ReadFile(subFS, "index.html")
			if err != nil {
				return nil, err
			}

			response = indexData
			return response, nil
		} else {
			return nil, err
		}
	}
	return response, nil
}

func (fs *FrontControllerService) V1FrontControllerServicePost(bodyJson map[string]interface{}) ([]byte, error) {
	return nil, nil
}

func (fs *FrontControllerService) V1FrontControllerServicePut(bodyJson map[string]interface{}) ([]byte, error) {
	return nil, nil
}

func (fs *FrontControllerService) V1FrontControllerServiceDelete(bodyJson map[string]interface{}) ([]byte, error) {
	return nil, nil
}
