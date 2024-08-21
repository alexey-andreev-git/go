package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"what-to.com/internal/config"
	"what-to.com/internal/repository"
)

type (
	FrontRoutesService struct {
		appRepository repository.Repository
		appConfig     *config.Config
		serviceFuncs  map[RequestType]ServiceFunc
	}
	RouteData struct {
		RouterLink string `json:"routerLink"`
		IconName   string `json:"iconName"`
		Title      string `json:"title"`
	}
	Route struct {
		Path       string     `json:"path"`
		Component  string     `json:"component,omitempty"`
		Data       *RouteData `json:"data,omitempty"`
		RedirectTo string     `json:"redirectTo,omitempty"`
		PathMatch  string     `json:"pathMatch,omitempty"`
	}
	Routes []Route
)

const (
	frontRoutesPath = "/front/routes"
)

const (
	FrontRoutesGet RequestType = iota
	FrontRoutesPost
	FrontRoutesPut
	FrontRoutesDelete
	FrontRoutesOptions
)

func NewFrontRoutesService(appConfig *config.Config, appRepo repository.Repository) *FrontRoutesService {
	s := &FrontRoutesService{
		appConfig:     appConfig,
		appRepository: appRepo,
	}
	s.registerServiceFuncs()
	return s
}

func setCorsHeader(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, authorization")
}

func (s *FrontRoutesService) registerServiceFuncs() {
	s.serviceFuncs = map[RequestType]ServiceFunc{
		FrontRoutesGet:     {s.V1FrontRoutesServiceGet, "GET", apiV1Path + frontRoutesPath + restWildcardPath},
		FrontRoutesPost:    {s.V1FrontRoutesServicePost, "POST", apiV1Path + frontRoutesPath + restWildcardPath},
		FrontRoutesPut:     {s.V1FrontRoutesServicePut, "PUT", apiV1Path + frontRoutesPath + restWildcardPath},
		FrontRoutesDelete:  {s.V1FrontRoutesServiceDelete, "DELETE", apiV1Path + frontRoutesPath + restWildcardPath},
		FrontRoutesOptions: {s.V1FrontRoutesServiceOptions, "OPTIONS", apiV1Path + frontRoutesPath + restWildcardPath},
	}
}

func (s *FrontRoutesService) GetServiceFuncs() map[RequestType]ServiceFunc {
	return s.serviceFuncs
}

func (s *FrontRoutesService) ServiceFunction(w http.ResponseWriter, r *http.Request, ver string, reqType RequestType) {
	if ver != "1" {
		ErrorHandler(s.appConfig.GetLogger(), w, errorMessage, fmt.Errorf("incorrect REST version"), http.StatusBadRequest)
		return
	}
	handlerFunc, ok := s.serviceFuncs[reqType]
	if !ok {
		ErrorHandler(s.appConfig.GetLogger(), w, errorMessage, fmt.Errorf("incorrect REST request type"), http.StatusBadRequest)
		return
	}
	bodyJson, err := GetRequestBodyJson(w, r)
	if err != nil {
		bodyJson = make(map[string]interface{})
	}
	if reqType == FrontRoutesOptions {
		setCorsHeader(w, r)
		w.WriteHeader(http.StatusOK)
		return
	}
	respJson, rerr := handlerFunc.Handler(bodyJson)
	if rerr != nil {
		ErrorHandler(s.appConfig.GetLogger(), w, errorMessage, rerr, http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	setCorsHeader(w, r)
	w.WriteHeader(http.StatusOK)
	w.Write(respJson)
}

func (s *FrontRoutesService) V1FrontRoutesServiceOptions(bodyJson map[string]interface{}) ([]byte, error) {
	return nil, nil
}

func (s *FrontRoutesService) V1FrontRoutesServiceGet(bodyJson map[string]interface{}) ([]byte, error) {
	routesMap := getDefaultRoutesMap()
	result, err := json.Marshal(routesMap)
	if err != nil {
		return nil, err
	}

	response := ([]byte)(result)
	return response, nil
}

func (s *FrontRoutesService) V1FrontRoutesServicePost(bodyJson map[string]interface{}) ([]byte, error) {
	response := ([]byte)(`{"message":"Fron Routes Post"}`)
	s.appConfig.GetLogger().Info("User created:" + string(response))
	return response, nil
}

func (s *FrontRoutesService) V1FrontRoutesServicePut(bodyJson map[string]interface{}) ([]byte, error) {
	rows, rerr := 0, fmt.Errorf("not implemented")
	return []byte(fmt.Sprintf(jsonOperationResultMsg, "updated", rows, rerr)), nil
}

func (s *FrontRoutesService) V1FrontRoutesServiceDelete(bodyJson map[string]interface{}) ([]byte, error) {
	rows, rerr := 0, fmt.Errorf("not implemented")
	return []byte(fmt.Sprintf(jsonOperationResultMsg, "deleted", rows, rerr)), nil
}

func getDefaultRoutesMap() Routes {
	return Routes{
		{
			Path:      "home",
			Component: "HomeComponent",
			Data:      &RouteData{RouterLink: "/home", IconName: "home", Title: "Home"},
		},
		{
			Path:      "register",
			Component: "RegisterComponent",
			Data:      &RouteData{RouterLink: "/register", IconName: "how_to_reg", Title: "Register"},
		},
		{
			Path:      "address",
			Component: "AddressFormComponent",
			Data:      &RouteData{RouterLink: "/address", IconName: "gite", Title: "Address"},
		},
		{
			Path:      "address1",
			Component: "AddressFormComponent",
			Data:      &RouteData{RouterLink: "/address1", IconName: "gite", Title: "Address"},
		},
		{
			Path:      "address2",
			Component: "AddressFormComponent",
			Data:      &RouteData{RouterLink: "/address2", IconName: "gite", Title: "Address"},
		},
		{
			Path:      "settings",
			Component: "AddressFormComponent",
			Data:      &RouteData{RouterLink: "/settings", IconName: "settings", Title: "Settings"},
		},
		{
			Path:       "",
			RedirectTo: "/home",
			PathMatch:  "full",
		},
		{
			Path:       "**",
			RedirectTo: "/home",
		},
	}
}
