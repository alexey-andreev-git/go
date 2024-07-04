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
	serviceFuncs  map[RequestType]func(bodyJson map[string]interface{}) ([]byte, error)
}

func NewEntityService(appConfig *config.Config, appRepo repository.Repository) *EntityService {
	es := &EntityService{
		appConfig:     appConfig,
		appRepository: appRepo,
	}
	es.registerServiceFuncs()
	return es
}

func (es *EntityService) registerServiceFuncs() {
	es.serviceFuncs = map[RequestType]func(bodyJson map[string]interface{}) ([]byte, error){
		EntityGet:           es.V1EntityServiceGet,
		EntityPost:          es.V1EntityServicePost,
		EntityPut:           es.V1EntityServicePut,
		EntityDelete:        es.V1EntityServiceDelete,
		EntityDataGet:       es.V1EntityDataServiceGet,
		EntityDataPost:      es.V1EntityDataServicePost,
		EntityDataPut:       es.V1EntityDataServicePut,
		EntityDataDelete:    es.V1EntityDataServiceDelete,
		EntityDataRefGet:    es.V1EntityDataRefServiceGet,
		EntityDataRefPost:   es.V1EntityDataRefServicePost,
		EntityDataRefPut:    es.V1EntityDataRefServicePut,
		EntityDataRefDelete: es.V1EntityDataRefServiceDelete,
	}
}

func (es *EntityService) GetServiceFunc(rt RequestType) func(bodyJson map[string]interface{}) ([]byte, error) {
	return es.serviceFuncs[rt]
}

func (es *EntityService) errorHandler(w http.ResponseWriter, message string, err error, status int) {
	es.appConfig.GetLogger().Error(message, err)
	http.Error(w, message+": "+err.Error(), status)
}

func (es *EntityService) EntityServiceFunction(w http.ResponseWriter, r *http.Request, ver string, reqType RequestType) {
	switch ver {
	case "1":
		if f, ok := es.serviceFuncs[reqType]; ok {
			bodyJson, err := es.V1EntityGetRequestBodyJson(w, r)
			if err != nil {
				return
			}
			respJson, reer := f(bodyJson)
			if reer == nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(respJson)
			} else {
				es.errorHandler(w, "Error processing request", reer, http.StatusBadRequest)
			}
		} else {
			es.errorHandler(w, "API error", fmt.Errorf("incorrect REST request type"), http.StatusBadRequest)
		}
	default:
		es.errorHandler(w, "API error", fmt.Errorf("incorrect REST version"), http.StatusBadRequest)
	}
}

func (es *EntityService) V1EntityGetRequestBodyJson(w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
	defer r.Body.Close()
	bodyJson := make(map[string]interface{})
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&bodyJson)
	if err != nil {
		es.errorHandler(w, "Error decoding JSON request body", err, http.StatusBadRequest)
		return nil, err
	}
	return bodyJson, nil
}

func (es *EntityService) V1EntityServiceGet(bodyJson map[string]interface{}) ([]byte, error) {
	entities, err := es.appRepository.GetEntity(bodyJson)
	if err != nil {
		return nil, err
	}
	response, jerr := json.Marshal(entities)
	if jerr != nil {
		return nil, jerr
	}
	return response, nil
}

func (es *EntityService) V1EntityServicePost(bodyJson map[string]interface{}) ([]byte, error) {
	entities, err := es.appRepository.CreateEntity(bodyJson)
	if err != nil {
		return nil, err
	}
	response, jerr := json.Marshal(entities)
	if jerr != nil {
		return nil, jerr
	}
	es.appConfig.GetLogger().Info("Entity created:" + string(response))
	return response, nil
}

func (es *EntityService) V1EntityServicePut(bodyJson map[string]interface{}) ([]byte, error) {
	result, err := es.appRepository.UpdateEntity(bodyJson)
	if err != nil {
		return nil, err
	}
	rows, rerr := result.RowsAffected()
	return []byte(fmt.Sprintf("{ \"message\": \"updated\", \"rows\": %v, \"error\": %v }", rows, rerr)), nil
}

func (es *EntityService) V1EntityServiceDelete(bodyJson map[string]interface{}) ([]byte, error) {
	result, err := es.appRepository.DeleteEntity(bodyJson)
	if err != nil {
		return nil, err
	}
	rows, rerr := result.RowsAffected()
	return []byte(fmt.Sprintf("{ \"message\": \"deleted\", \"rows\": %v, \"error\": %v }", rows, rerr)), nil
}

func (es *EntityService) V1EntityDataServiceGet(bodyJson map[string]interface{}) ([]byte, error) {
	entityData, err := es.appRepository.GetEntityData(bodyJson)
	if err != nil {
		return nil, err
	}
	response, jerr := json.Marshal(entityData)
	if jerr != nil {
		return nil, jerr
	}
	return response, nil
}

func (es *EntityService) V1EntityDataServicePost(bodyJson map[string]interface{}) ([]byte, error) {
	entityData, err := es.appRepository.CreateEntityData(bodyJson)
	if err != nil {
		return nil, err
	}
	response, jerr := json.Marshal(entityData)
	if jerr != nil {
		return nil, jerr
	}
	es.appConfig.GetLogger().Info("Entity data created:" + string(response))
	return response, nil
}

func (es *EntityService) V1EntityDataServicePut(bodyJson map[string]interface{}) ([]byte, error) {
	result, err := es.appRepository.UpdateEntityData(bodyJson)
	if err != nil {
		return nil, err
	}
	rows, rerr := result.RowsAffected()
	return []byte(fmt.Sprintf("{ \"message\": \"updated\", \"rows\": %v, \"error\": %v }", rows, rerr)), nil
}

func (es *EntityService) V1EntityDataServiceDelete(bodyJson map[string]interface{}) ([]byte, error) {
	result, err := es.appRepository.DeleteEntityData(bodyJson)
	if err != nil {
		return nil, err
	}
	rows, rerr := result.RowsAffected()
	return []byte(fmt.Sprintf("{ \"message\": \"deleted\", \"rows\": %v, \"error\": %v }", rows, rerr)), nil
}

func (es *EntityService) V1EntityDataRefServiceGet(bodyJson map[string]interface{}) ([]byte, error) {
	entityDataRef, err := es.appRepository.GetEntityDataRef(bodyJson)
	if err != nil {
		return nil, err
	}
	response, jerr := json.Marshal(entityDataRef)
	if jerr != nil {
		return nil, jerr
	}
	return response, nil
}

func (es *EntityService) V1EntityDataRefServicePost(bodyJson map[string]interface{}) ([]byte, error) {
	entityDataRef, err := es.appRepository.CreateEntityDataRef(bodyJson)
	if err != nil {
		return nil, err
	}
	response, jerr := json.Marshal(entityDataRef)
	if jerr != nil {
		return nil, jerr
	}
	es.appConfig.GetLogger().Info("Entity data reference created:" + string(response))
	return response, nil
}

func (es *EntityService) V1EntityDataRefServicePut(bodyJson map[string]interface{}) ([]byte, error) {
	result, err := es.appRepository.UpdateEntityDataRef(bodyJson)
	if err != nil {
		return nil, err
	}
	rows, rerr := result.RowsAffected()
	return []byte(fmt.Sprintf("{ \"message\": \"updated\", \"rows\": %v, \"error\": %v }", rows, rerr)), nil
}

func (es *EntityService) V1EntityDataRefServiceDelete(bodyJson map[string]interface{}) ([]byte, error) {
	result, err := es.appRepository.DeleteEntityDataRef(bodyJson)
	if err != nil {
		return nil, err
	}
	rows, rerr := result.RowsAffected()
	return []byte(fmt.Sprintf("{ \"message\": \"deleted\", \"rows\": %v, \"error\": %v }", rows, rerr)), nil
}
