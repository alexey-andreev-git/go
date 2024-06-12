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
}

func NewEntityService(appConfig *config.Config, appRepo repository.Repository) *EntityService {
	return &EntityService{appConfig: appConfig, appRepository: appRepo}
}

func (es *EntityService) errorHandler(w http.ResponseWriter, message string, err error, status int) {
	es.appConfig.GetLogger().Error(message, err)
	http.Error(w, message+": "+err.Error(), status)
}

func (es *EntityService) EntityServiceFunctionGet(w http.ResponseWriter, r *http.Request, ver string) {
	switch ver {
	case "1":
		es.V1EntityServiceGet(w, r)
	default:
		http.Error(w, "incorrect REST version", http.StatusBadRequest)
	}
}

func (es *EntityService) EntityServiceFunctionPost(w http.ResponseWriter, r *http.Request, ver string) {
	switch ver {
	case "1":
		es.V1EntityServicePost(w, r)
	default:
		http.Error(w, "incorrect REST version", http.StatusBadRequest)
	}
}

func (es *EntityService) EntityServiceFunctionPut(w http.ResponseWriter, r *http.Request, ver string) {
	switch ver {
	case "1":
		es.V1EntityServicePut(w, r)
	default:
		http.Error(w, "incorrect REST version", http.StatusBadRequest)
	}
}

func (es *EntityService) EntityServiceFunctionDelete(w http.ResponseWriter, r *http.Request, ver string) {
	switch ver {
	case "1":
		es.V1EntityServiceDelete(w, r)
	default:
		http.Error(w, "incorrect REST version", http.StatusBadRequest)
	}
}

func (es *EntityService) EntityDataServiceFunctionGet(w http.ResponseWriter, r *http.Request, ver string) {
	switch ver {
	case "1":
		es.V1EntityServiceGetEntityData(w, r)
	default:
		http.Error(w, "incorrect REST version", http.StatusBadRequest)
	}
}

func (es *EntityService) EntityDataServiceFunctionPost(w http.ResponseWriter, r *http.Request, ver string) {
	switch ver {
	case "1":
		es.V1EntityServicePostEntityData(w, r)
	default:
		http.Error(w, "incorrect REST version", http.StatusBadRequest)
	}
}

func (es *EntityService) EntityDataServiceFunctionPut(w http.ResponseWriter, r *http.Request, ver string) {
	switch ver {
	case "1":
		es.V1EntityServicePutEntityData(w, r)
	default:
		http.Error(w, "incorrect REST version", http.StatusBadRequest)
	}
}

func (es *EntityService) EntityDataServiceFunctionDelete(w http.ResponseWriter, r *http.Request, ver string) {
	switch ver {
	case "1":
		es.V1EntityServiceDeleteEntityData(w, r)
	default:
		http.Error(w, "incorrect REST version", http.StatusBadRequest)
	}
}

func (es *EntityService) EntityDataRefServiceFunctionGet(w http.ResponseWriter, r *http.Request, ver string) {
	switch ver {
	case "1":
		es.V1EntityServiceGetEntityDataRef(w, r)
	default:
		http.Error(w, "incorrect REST version", http.StatusBadRequest)
	}
}

func (es *EntityService) EntityDataRefServiceFunctionPost(w http.ResponseWriter, r *http.Request, ver string) {
	switch ver {
	case "1":
		es.V1EntityServicePostEntityDataRef(w, r)
	default:
		http.Error(w, "incorrect REST version", http.StatusBadRequest)
	}
}

func (es *EntityService) EntityDataRefServiceFunctionPut(w http.ResponseWriter, r *http.Request, ver string) {
	switch ver {
	case "1":
		es.V1EntityServicePutEntityDataRef(w, r)
	default:
		http.Error(w, "incorrect REST version", http.StatusBadRequest)
	}
}

func (es *EntityService) EntityDataRefServiceFunctionDelete(w http.ResponseWriter, r *http.Request, ver string) {
	switch ver {
	case "1":
		es.V1EntityServiceDeleteEntityDataRef(w, r)
	default:
		http.Error(w, "incorrect REST version", http.StatusBadRequest)
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

func (es *EntityService) V1EntityServiceGet(w http.ResponseWriter, r *http.Request) {
	bodyJson, err := es.V1EntityGetRequestBodyJson(w, r)
	if err != nil {
		return
	}
	entities, eerr := es.appRepository.GetEntity(bodyJson)
	if eerr != nil {
		es.errorHandler(w, "Error getting entity", eerr, http.StatusNotFound)
		return
	}
	// Ensure Content-Type is set for all responses
	w.Header().Set("Content-Type", "application/json")
	response, jerr := json.Marshal(entities)
	if jerr != nil {
		es.errorHandler(w, "Error marshalling JSON response", jerr, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (es *EntityService) V1EntityServicePost(w http.ResponseWriter, r *http.Request) {
	bodyJson, err := es.V1EntityGetRequestBodyJson(w, r)
	if err != nil {
		return
	}
	entities, eerr := es.appRepository.CreateEntity(bodyJson)
	if eerr != nil {
		es.errorHandler(w, "Error creating entity", eerr, http.StatusBadRequest)
		return
	}
	// Ensure Content-Type is set for all responses
	w.Header().Set("Content-Type", "application/json")
	response, jerr := json.Marshal(entities)
	if jerr != nil {
		es.errorHandler(w, "Error marshalling JSON response", jerr, http.StatusBadRequest)
		return
	}
	es.appConfig.GetLogger().Info("Entity created:" + string(response))
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (es *EntityService) V1EntityServicePut(w http.ResponseWriter, r *http.Request) {
	bodyJson, err := es.V1EntityGetRequestBodyJson(w, r)
	if err != nil {
		return
	}

	result, uerr := es.appRepository.UpdateEntity(bodyJson)
	if uerr != nil {
		es.errorHandler(w, "Error updating entity", uerr, http.StatusBadRequest)
		return
	}
	rows, rerr := result.RowsAffected()
	// Ensure Content-Type is set for all responses
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("{ \"message\": \"updated\", \"rows\": %v, \"error\": %v }", rows, rerr)))
}

func (es *EntityService) V1EntityServiceDelete(w http.ResponseWriter, r *http.Request) {
	bodyJson, err := es.V1EntityGetRequestBodyJson(w, r)
	if err != nil {
		return
	}

	result, derr := es.appRepository.DeleteEntity(bodyJson)
	if derr != nil {
		es.errorHandler(w, "Error deleting entity", derr, http.StatusBadRequest)
		return
	}
	rows, rerr := result.RowsAffected()
	// Ensure Content-Type is set for all responses
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("{ \"message\": \"deleted\", \"rows\": %v, \"error\": %v }", rows, rerr)))
}

// Add these methods

func (es *EntityService) V1EntityServiceGetEntityData(w http.ResponseWriter, r *http.Request) {
	bodyJson, err := es.V1EntityGetRequestBodyJson(w, r)
	if err != nil {
		return
	}
	entityData, eerr := es.appRepository.GetEntityData(bodyJson)
	if eerr != nil {
		es.errorHandler(w, "Error getting entity data", eerr, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response, jerr := json.Marshal(entityData)
	if jerr != nil {
		es.errorHandler(w, "Error marshalling JSON response", jerr, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (es *EntityService) V1EntityServicePostEntityData(w http.ResponseWriter, r *http.Request) {
	bodyJson, err := es.V1EntityGetRequestBodyJson(w, r)
	if err != nil {
		return
	}
	entityData, eerr := es.appRepository.CreateEntityData(bodyJson)
	if eerr != nil {
		es.errorHandler(w, "Error creating entity data", eerr, http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response, jerr := json.Marshal(entityData)
	if jerr != nil {
		es.errorHandler(w, "Error marshalling JSON response", jerr, http.StatusBadRequest)
		return
	}
	es.appConfig.GetLogger().Info("Entity data created:" + string(response))
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (es *EntityService) V1EntityServicePutEntityData(w http.ResponseWriter, r *http.Request) {
	bodyJson, err := es.V1EntityGetRequestBodyJson(w, r)
	if err != nil {
		return
	}

	result, uerr := es.appRepository.UpdateEntityData(bodyJson)
	if uerr != nil {
		es.errorHandler(w, "Error updating entity data", uerr, http.StatusBadRequest)
		return
	}
	rows, rerr := result.RowsAffected()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("{ \"message\": \"updated\", \"rows\": %v, \"error\": %v }", rows, rerr)))
}

func (es *EntityService) V1EntityServiceDeleteEntityData(w http.ResponseWriter, r *http.Request) {
	bodyJson, err := es.V1EntityGetRequestBodyJson(w, r)
	if err != nil {
		return
	}

	result, derr := es.appRepository.DeleteEntityData(bodyJson)
	if derr != nil {
		es.errorHandler(w, "Error deleting entity data", derr, http.StatusBadRequest)
		return
	}
	rows, rerr := result.RowsAffected()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("{ \"message\": \"deleted\", \"rows\": %v, \"error\": %v }", rows, rerr)))
}

func (es *EntityService) V1EntityServiceGetEntityDataRef(w http.ResponseWriter, r *http.Request) {
	bodyJson, err := es.V1EntityGetRequestBodyJson(w, r)
	if err != nil {
		return
	}
	entityDataRef, eerr := es.appRepository.GetEntityDataRef(bodyJson)
	if eerr != nil {
		es.errorHandler(w, "Error getting entity data reference", eerr, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response, jerr := json.Marshal(entityDataRef)
	if jerr != nil {
		es.errorHandler(w, "Error marshalling JSON response", jerr, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (es *EntityService) V1EntityServicePostEntityDataRef(w http.ResponseWriter, r *http.Request) {
	bodyJson, err := es.V1EntityGetRequestBodyJson(w, r)
	if err != nil {
		return
	}
	entityDataRef, eerr := es.appRepository.CreateEntityDataRef(bodyJson)
	if eerr != nil {
		es.errorHandler(w, "Error creating entity data reference", eerr, http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response, jerr := json.Marshal(entityDataRef)
	if jerr != nil {
		es.errorHandler(w, "Error marshalling JSON response", jerr, http.StatusBadRequest)
		return
	}
	es.appConfig.GetLogger().Info("Entity data reference created:" + string(response))
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (es *EntityService) V1EntityServicePutEntityDataRef(w http.ResponseWriter, r *http.Request) {
	bodyJson, err := es.V1EntityGetRequestBodyJson(w, r)
	if err != nil {
		return
	}

	result, uerr := es.appRepository.UpdateEntityDataRef(bodyJson)
	if uerr != nil {
		es.errorHandler(w, "Error updating entity data reference", uerr, http.StatusBadRequest)
		return
	}
	rows, rerr := result.RowsAffected()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("{ \"message\": \"updated\", \"rows\": %v, \"error\": %v }", rows, rerr)))
}

func (es *EntityService) V1EntityServiceDeleteEntityDataRef(w http.ResponseWriter, r *http.Request) {
	bodyJson, err := es.V1EntityGetRequestBodyJson(w, r)
	if err != nil {
		return
	}

	result, derr := es.appRepository.DeleteEntityDataRef(bodyJson)
	if derr != nil {
		es.errorHandler(w, "Error deleting entity data reference", derr, http.StatusBadRequest)
		return
	}
	rows, rerr := result.RowsAffected()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("{ \"message\": \"deleted\", \"rows\": %v, \"error\": %v }", rows, rerr)))
}
