package service

import (
	"encoding/json"
	"net/http"
	"strconv"

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

func (es *EntityService) V1EntityGetRequestBodyJson(w http.ResponseWriter, r *http.Request) (map[string]string, error) {
	defer r.Body.Close()
	bodyJson := make(map[string]string)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&bodyJson)
	if err != nil {
		es.errorHandler(w, "Error decoding JSON request body", err, http.StatusServiceUnavailable)
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
		es.errorHandler(w, "Error getting entity", eerr, http.StatusServiceUnavailable)
		return
	}
	// Ensure Content-Type is set for all responses
	w.Header().Set("Content-Type", "application/json")
	response, jerr := json.Marshal(entities)
	if jerr != nil {
		es.errorHandler(w, "Error marshalling JSON response", jerr, http.StatusServiceUnavailable)
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
		es.errorHandler(w, "Error creating entity", eerr, http.StatusServiceUnavailable)
		return
	}
	// Ensure Content-Type is set for all responses
	w.Header().Set("Content-Type", "application/json")
	response, jerr := json.Marshal(entities)
	if jerr != nil {
		es.errorHandler(w, "Error marshalling JSON response", jerr, http.StatusServiceUnavailable)
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

	idStr, ok := bodyJson["Id"]
	if !ok {
		http.Error(w, "Missing entity ID in request body", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid entity ID", http.StatusBadRequest)
		return
	}

	delete(bodyJson, "Id") // Remove ID from bodyJson to avoid updating it
	err = es.appRepository.UpdateEntity(id, bodyJson)
	if err != nil {
		es.errorHandler(w, "Error updating entity", err, http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (es *EntityService) V1EntityServiceDelete(w http.ResponseWriter, r *http.Request) {
	bodyJson, err := es.V1EntityGetRequestBodyJson(w, r)
	if err != nil {
		return
	}

	idStr, ok := bodyJson["Id"]
	if !ok {
		http.Error(w, "Missing entity ID in request body", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid entity ID", http.StatusBadRequest)
		return
	}

	err = es.appRepository.DeleteEntity(id)
	if err != nil {
		es.errorHandler(w, "Error deleting entity", err, http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
