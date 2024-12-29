package service

import (
	"fmt"
	"net/http"

	"what-to.com/internal/config"
	"what-to.com/internal/models"
	"what-to.com/internal/repository"
)

type (
	FrontService struct {
		appRepository repository.Repository
		appConfig     *config.Config
		serviceFuncs  map[RequestType]ServiceFunc
	}
)

const (
	frontRoutesPath = "/front/routes"
	frontFormsPath  = "/front/forms"
)

const (
	FrontRoutesGet RequestType = iota
	FrontFormDtaGet
	FrontRoutesPost
	FrontRoutesPut
	FrontRoutesDelete
)

func NewFrontService(appConfig *config.Config, appRepo repository.Repository) *FrontService {
	s := &FrontService{
		appConfig:     appConfig,
		appRepository: appRepo,
	}
	s.registerServiceFuncs()
	return s
}

func (s *FrontService) registerServiceFuncs() {
	s.serviceFuncs = map[RequestType]ServiceFunc{
		FrontRoutesGet:    {s.V1FrontServiceRoutesGet, "GET", apiV1Path + frontRoutesPath + restWildcardPath},
		FrontFormDtaGet:   {s.V1FrontServiceFormDataGet, "GET", apiV1Path + frontRoutesPath + restWildcardPath},
		FrontRoutesPost:   {s.V1FrontServicePost, "POST", apiV1Path + frontRoutesPath + restWildcardPath},
		FrontRoutesPut:    {s.V1FrontServicePut, "PUT", apiV1Path + frontRoutesPath + restWildcardPath},
		FrontRoutesDelete: {s.V1FrontServiceDelete, "DELETE", apiV1Path + frontRoutesPath + restWildcardPath},
	}
}

func (s *FrontService) GetServiceFuncs() map[RequestType]ServiceFunc {
	return s.serviceFuncs
}

func (s *FrontService) ServiceFunction(w http.ResponseWriter, r *http.Request, ver string, reqType RequestType) {
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
	respJson, rerr := handlerFunc.Handler(bodyJson)
	if rerr != nil {
		ErrorHandler(s.appConfig.GetLogger(), w, errorMessage, rerr, http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respJson)
}

func (s *FrontService) V1FrontServiceRoutesGet(data interface{}) ([]byte, error) {
	routes, ok := data.(*models.Routes)
	if !ok {
		err := fmt.Errorf("incorrect data type")
		s.appConfig.GetLogger().Error("Fron Routes Get: incorrect data type", err)
		return nil, err
	}
	*routes = getDefaultRoutesMap()
	return nil, nil
}

func (s *FrontService) V1FrontServicePost(data interface{}) ([]byte, error) {
	d, ok := data.(map[string]interface{})
	if !ok {
		err := fmt.Errorf("incorrect data type")
		s.appConfig.GetLogger().Error("Fron Routes Post: incorrect data type", err)
		return nil, err
	}
	d["message"] = "Fron Routes Post"
	s.appConfig.GetLogger().Info("Fron Routes Post: " + d["message"].(string))
	return nil, nil
}

func (s *FrontService) V1FrontServicePut(bodyJson interface{}) ([]byte, error) {
	rows, rerr := 0, fmt.Errorf("not implemented")
	return []byte(fmt.Sprintf(jsonOperationResultMsg, "updated", rows, rerr)), nil
}

func (s *FrontService) V1FrontServiceDelete(bodyJson interface{}) ([]byte, error) {
	rows, rerr := 0, fmt.Errorf("not implemented")
	return []byte(fmt.Sprintf(jsonOperationResultMsg, "deleted", rows, rerr)), nil
}

func (s *FrontService) V1FrontServiceFormDataGet(data interface{}) ([]byte, error) {
	routes, ok := data.(*models.Controls)
	if !ok {
		err := fmt.Errorf("incorrect data type")
		s.appConfig.GetLogger().Error("Fron Routes Get: incorrect data type", err)
		return nil, err
	}
	*routes = getDefaultFormsMap()
	return nil, nil
}

func getDefaultRoutesMap() models.Routes {
	return models.Routes{
		{
			Path:      "home",
			Component: "HomeComponent",
			Data:      &models.RouteData{RouterLink: "/home", IconName: "home", Title: "Home"},
		},
		{
			Path:      "register",
			Component: "RegisterComponent",
			Data:      &models.RouteData{RouterLink: "/register", IconName: "how_to_reg", Title: "Register"},
		},
		{
			Path:      "address",
			Component: "AddressFormComponent",
			Data:      &models.RouteData{RouterLink: "/address", IconName: "gite", Title: "Address"},
		},
		{
			Path:      "address1",
			Component: "AddressFormComponent",
			Data:      &models.RouteData{RouterLink: "/address1", IconName: "gite", Title: "Address"},
		},
		{
			Path:      "address2",
			Component: "AddressFormComponent",
			Data:      &models.RouteData{RouterLink: "/address2", IconName: "gite", Title: "Address"},
		},
		{
			Path:      "settings",
			Component: "AddressFormComponent",
			Data:      &models.RouteData{RouterLink: "/settings", IconName: "settings", Title: "Settings"},
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

func getDefaultFormsMap() models.Controls {
	return models.Controls{
		Controls: []models.ControlData{
			{Name: "firstName", Type: "text", Placeholder: "First Name", Validators: []string{"required"}},
			{Name: "lastName", Type: "text", Placeholder: "Last Name", Validators: []string{"required"}},
			{Name: "email", Type: "email", Placeholder: "Email", Validators: []string{"required", "email"}},
			{Name: "password", Type: "password", Placeholder: "Password", Validators: []string{"required", "minLength:6"}},
			{Name: "confirmPassword", Type: "password", Placeholder: "Confirm Password", Validators: []string{"required", "minLength:6"}},
			{Name: "test", Type: "submit", Placeholder: "Test", Validators: nil},
		},
	}
}
