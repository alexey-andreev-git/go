package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"what-to.com/internal/config"
	"what-to.com/internal/repository"
)

// Struct of entity service
type EntityService struct {
	appRepository repository.Repository
	appConfig     *config.Config
	serviceFuncs  map[RequestType]ServiceFunc
}

const (
	EntityGet RequestType = iota
	EntityPost
	EntityPut
	EntityDelete
	EntityDataGet
	EntityDataPost
	EntityDataPut
	EntityDataDelete
	EntityDataRefGet
	EntityDataRefPost
	EntityDataRefPut
	EntityDataRefDelete
)

const (
	entityPath             = "/entity"
	entityDataPath         = "/entities_data"
	entityDataRefPath      = "/entities_data_reference"
	jsonOperationResultMsg = "{ \"message\": \"%s\", \"rows\": %v, \"error\": %v }"
)

func NewEntityService(appConfig *config.Config, appRepo repository.Repository) *EntityService {
	es := &EntityService{
		appConfig:     appConfig,
		appRepository: appRepo,
	}
	es.registerServiceFuncs()
	return es
}

func (s *EntityService) registerServiceFuncs() {
	s.serviceFuncs = map[RequestType]ServiceFunc{
		EntityGet:           {s.V1EntityServiceGet, "GET", apiV1Path + entityPath + restWildcardPath},
		EntityPost:          {s.V1EntityServicePost, "POST", apiV1Path + entityPath + restWildcardPath},
		EntityPut:           {s.V1EntityServicePut, "PUT", apiV1Path + entityPath + restWildcardPath},
		EntityDelete:        {s.V1EntityServiceDelete, "DELETE", apiV1Path + entityPath + restWildcardPath},
		EntityDataGet:       {s.V1EntityDataServiceGet, "GET", apiV1Path + entityDataPath + restWildcardPath},
		EntityDataPost:      {s.V1EntityDataServicePost, "POST", apiV1Path + entityDataPath + restWildcardPath},
		EntityDataPut:       {s.V1EntityDataServicePut, "PUT", apiV1Path + entityDataPath + restWildcardPath},
		EntityDataDelete:    {s.V1EntityDataServiceDelete, "DELETE", apiV1Path + entityDataPath + restWildcardPath},
		EntityDataRefGet:    {s.V1EntityDataRefServiceGet, "GET", apiV1Path + entityDataRefPath + restWildcardPath},
		EntityDataRefPost:   {s.V1EntityDataRefServicePost, "POST", apiV1Path + entityDataRefPath + restWildcardPath},
		EntityDataRefPut:    {s.V1EntityDataRefServicePut, "PUT", apiV1Path + entityDataRefPath + restWildcardPath},
		EntityDataRefDelete: {s.V1EntityDataRefServiceDelete, "DELETE", apiV1Path + entityDataRefPath + restWildcardPath},
	}
}

func (s *EntityService) GetServiceFuncs() map[RequestType]ServiceFunc {
	return s.serviceFuncs
}

func (s *EntityService) ServiceFunction(w http.ResponseWriter, r *http.Request, ver string, reqType RequestType) {
	if ver != "1" {
		ErrorHandler(s.appConfig.GetLogger(), w, errorMessage, fmt.Errorf("incorrect REST version"), http.StatusBadRequest)
		return
	}
	handlerFunc, ok := s.serviceFuncs[reqType]
	if !ok {
		ErrorHandler(s.appConfig.GetLogger(), w, errorMessage, fmt.Errorf("incorrect REST type"), http.StatusBadRequest)
		return
	}
	// bodyJson, err := s.V1EntityGetRequestBodyJson(w, r)
	bodyJson, err := GetRequestBodyJson(w, r)
	if err != nil {
		ErrorHandler(s.appConfig.GetLogger(), w, "Error decoding JSON request body", err, http.StatusBadRequest)
		return
	}
	respJson, reer := handlerFunc.Handler(bodyJson)
	if reer == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(respJson)
	} else {
		ErrorHandler(s.appConfig.GetLogger(), w, "Error processing request", reer, http.StatusBadRequest)
	}
}

func (s *EntityService) V1EntityServiceGet(bodyJson map[string]interface{}) ([]byte, error) {
	entities, err := s.appRepository.(*repository.PgRepository).GetEntity(bodyJson)
	if err != nil {
		return nil, err
	}
	response, jerr := json.Marshal(entities)
	if jerr != nil {
		return nil, jerr
	}
	return response, nil
}

func (s *EntityService) V1EntityServicePost(bodyJson map[string]interface{}) ([]byte, error) {
	entities, err := s.appRepository.(*repository.PgRepository).CreateEntity(bodyJson)
	if err != nil {
		return nil, err
	}
	response, jerr := json.Marshal(entities)
	if jerr != nil {
		return nil, jerr
	}
	s.appConfig.GetLogger().Info("Entity created:" + string(response))
	return response, nil
}

func (s *EntityService) V1EntityServicePut(bodyJson map[string]interface{}) ([]byte, error) {
	result, err := s.appRepository.(*repository.PgRepository).UpdateEntity(bodyJson)
	if err != nil {
		return nil, err
	}
	rows, rerr := result.RowsAffected()
	return []byte(fmt.Sprintf(jsonOperationResultMsg, "updated", rows, rerr)), nil
}

func (s *EntityService) V1EntityServiceDelete(bodyJson map[string]interface{}) ([]byte, error) {
	result, err := s.appRepository.(*repository.PgRepository).DeleteEntity(bodyJson)
	if err != nil {
		return nil, err
	}
	rows, rerr := result.RowsAffected()
	return []byte(fmt.Sprintf(jsonOperationResultMsg, "deleted", rows, rerr)), nil
}

func (s *EntityService) V1EntityDataServiceGet(bodyJson map[string]interface{}) ([]byte, error) {
	entityData, err := s.appRepository.(*repository.PgRepository).GetEntityData(bodyJson)
	if err != nil {
		return nil, err
	}
	response, jerr := json.Marshal(entityData)
	if jerr != nil {
		return nil, jerr
	}
	return response, nil
}

func (s *EntityService) V1EntityDataServicePost(bodyJson map[string]interface{}) ([]byte, error) {
	entityData, err := s.appRepository.(*repository.PgRepository).CreateEntityData(bodyJson)
	if err != nil {
		return nil, err
	}
	response, jerr := json.Marshal(entityData)
	if jerr != nil {
		return nil, jerr
	}
	s.appConfig.GetLogger().Info("Entity data created:" + string(response))
	return response, nil
}

func (s *EntityService) V1EntityDataServicePut(bodyJson map[string]interface{}) ([]byte, error) {
	result, err := s.appRepository.(*repository.PgRepository).UpdateEntityData(bodyJson)
	if err != nil {
		return nil, err
	}
	rows, rerr := result.RowsAffected()
	return []byte(fmt.Sprintf(jsonOperationResultMsg, "updated", rows, rerr)), nil
}

func (s *EntityService) V1EntityDataServiceDelete(bodyJson map[string]interface{}) ([]byte, error) {
	result, err := s.appRepository.(*repository.PgRepository).DeleteEntityData(bodyJson)
	if err != nil {
		return nil, err
	}
	rows, rerr := result.RowsAffected()
	return []byte(fmt.Sprintf(jsonOperationResultMsg, "deleted", rows, rerr)), nil
}

func (s *EntityService) V1EntityDataRefServiceGet(bodyJson map[string]interface{}) ([]byte, error) {
	entityDataRef, err := s.appRepository.(*repository.PgRepository).GetEntityDataRef(bodyJson)
	if err != nil {
		return nil, err
	}
	response, jerr := json.Marshal(entityDataRef)
	if jerr != nil {
		return nil, jerr
	}
	return response, nil
}

func (s *EntityService) V1EntityDataRefServicePost(bodyJson map[string]interface{}) ([]byte, error) {
	entityDataRef, err := s.appRepository.(*repository.PgRepository).CreateEntityDataRef(bodyJson)
	if err != nil {
		return nil, err
	}
	response, jerr := json.Marshal(entityDataRef)
	if jerr != nil {
		return nil, jerr
	}
	s.appConfig.GetLogger().Info("Entity data reference created:" + string(response))
	return response, nil
}

func (s *EntityService) V1EntityDataRefServicePut(bodyJson map[string]interface{}) ([]byte, error) {
	result, err := s.appRepository.(*repository.PgRepository).UpdateEntityDataRef(bodyJson)
	if err != nil {
		return nil, err
	}
	rows, rerr := result.RowsAffected()
	return []byte(fmt.Sprintf(jsonOperationResultMsg, "updated", rows, rerr)), nil
}

func (s *EntityService) V1EntityDataRefServiceDelete(bodyJson map[string]interface{}) ([]byte, error) {
	result, err := s.appRepository.(*repository.PgRepository).DeleteEntityDataRef(bodyJson)
	if err != nil {
		return nil, err
	}
	rows, rerr := result.RowsAffected()
	return []byte(fmt.Sprintf(jsonOperationResultMsg, "deleted", rows, rerr)), nil
}
