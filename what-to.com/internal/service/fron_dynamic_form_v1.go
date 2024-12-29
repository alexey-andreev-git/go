package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"what-to.com/internal/config"
	"what-to.com/internal/repository"
)

type (
	DynamicFormService struct {
		appRepository repository.Repository
		appConfig     *config.Config
		serviceFuncs  map[RequestType]ServiceFunc
	}
)

const (
	dynamicFormPath = "/form/data"
)

const (
	DynamicFormGet RequestType = iota
	DynamicFormPost
	DynamicFormPut
	DynamicFormDelete
)

func NewDynamicFormService(appConfig *config.Config, appRepo repository.Repository) *DynamicFormService {
	s := &DynamicFormService{
		appConfig:     appConfig,
		appRepository: appRepo,
	}
	s.registerServiceFuncs()
	return s
}

func (s *DynamicFormService) registerServiceFuncs() {
	s.serviceFuncs = map[RequestType]ServiceFunc{
		DynamicFormGet:    {s.V1DynamicFormServiceGet, "GET", apiV1Path + dynamicFormPath + restWildcardPath},
		DynamicFormPost:   {s.V1DynamicFormServicePost, "POST", apiV1Path + dynamicFormPath + restWildcardPath},
		DynamicFormPut:    {s.V1DynamicFormServicePut, "PUT", apiV1Path + dynamicFormPath + restWildcardPath},
		DynamicFormDelete: {s.V1DynamicFormServiceDelete, "DELETE", apiV1Path + dynamicFormPath + restWildcardPath},
	}
}

func (s *DynamicFormService) GetServiceFuncs() map[RequestType]ServiceFunc {
	return s.serviceFuncs
}

func (s *DynamicFormService) ServiceFunction(w http.ResponseWriter, r *http.Request, ver string, reqType RequestType) {
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

func (s *DynamicFormService) V1DynamicFormServiceGet(bodyJson interface{}) ([]byte, error) {
	routesMap := getDefaultFormsMap()
	result, err := json.Marshal(routesMap)
	if err != nil {
		return nil, err
	}

	response := ([]byte)(result)
	return response, nil
}

func (s *DynamicFormService) V1DynamicFormServicePost(bodyJson interface{}) ([]byte, error) {
	response := ([]byte)(`{"message":"Fron Routes Post"}`)
	s.appConfig.GetLogger().Info("User created:" + string(response))
	return response, nil
}

func (s *DynamicFormService) V1DynamicFormServicePut(bodyJson interface{}) ([]byte, error) {
	rows, rerr := 0, fmt.Errorf("not implemented")
	return []byte(fmt.Sprintf(jsonOperationResultMsg, "updated", rows, rerr)), nil
}

func (s *DynamicFormService) V1DynamicFormServiceDelete(bodyJson interface{}) ([]byte, error) {
	rows, rerr := 0, fmt.Errorf("not implemented")
	return []byte(fmt.Sprintf(jsonOperationResultMsg, "deleted", rows, rerr)), nil
}

// func getDefaultFormsMap() Controls {
// 	return Controls{
// 		Controls: []ControlData{
// 			{Name: "firstName", Type: "text", Placeholder: "First Name", Validators: []string{"required"}},
// 			{Name: "lastName", Type: "text", Placeholder: "Last Name", Validators: []string{"required"}},
// 			{Name: "email", Type: "email", Placeholder: "Email", Validators: []string{"required", "email"}},
// 			{Name: "password", Type: "password", Placeholder: "Password", Validators: []string{"required", "minLength:6"}},
// 			{Name: "confirmPassword", Type: "password", Placeholder: "Confirm Password", Validators: []string{"required", "minLength:6"}},
// 			{Name: "test", Type: "submit", Placeholder: "Test", Validators: nil},
// 		},
// 	}
// }
