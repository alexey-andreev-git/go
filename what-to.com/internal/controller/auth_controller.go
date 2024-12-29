package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"what-to.com/internal/config"
	"what-to.com/internal/models"
	"what-to.com/internal/service"
)

type AuthControllerV1 struct {
	httpHandlers HttpHandlersT
	config       *config.Config
	newService   func() service.Service
}

const (
	authPath = "/auth"
)

func NewAuthControllerV1(appConfig *config.Config, newSvc func() service.Service) *AuthControllerV1 {
	c := &AuthControllerV1{
		config:     appConfig,
		newService: newSvc,
	}
	c.httpHandlers = HttpHandlersT{
		{Method: "GET", Handler: c.AuthGet, Path: apiV1Path + authPath + restWildcardPath},
		{Method: "POST", Handler: c.AuthPost, Path: apiV1Path + authPath + restWildcardPath},
	}
	return c
}

func (c *AuthControllerV1) AuthGet(w http.ResponseWriter, r *http.Request) {
	// Read the token from the header
	token := r.Header.Get("Authorization")
	// Check if the token is in the correct format
	if token == "" {
		ErrorHandler(c.config.GetLogger(), w, errorMessage, fmt.Errorf("no token"), http.StatusForbidden)
		return
	}
	bearer := "bearer"
	tokenLen := len(bearer)
	if strings.ToLower(token[:tokenLen]) != bearer {
		ErrorHandler(c.config.GetLogger(), w, errorMessage, fmt.Errorf("incorrect token"), http.StatusBadRequest)
		return
	}
	// Remove any leading spaces from the token
	for tokenLen < len(token) && token[tokenLen] == ' ' {
		tokenLen++
	}
	// Remove "bearer " from the token
	token = token[tokenLen:]
	// Create the bodyJson map with the token
	userData := &struct {
		User *models.User `json:"user"`
	}{User: &models.User{Token: token}}
	// Call the service function
	svc := c.newService().(*service.AuthService)
	f := svc.GetServiceFuncs()[service.AuthGet]

	if _, err := f.Handler(userData.User); err != nil {
		ErrorHandler(c.config.GetLogger(), w, errorMessage, err, http.StatusBadRequest)
		return
	}
	SetResponseJson(w, userData)
}

func (c *AuthControllerV1) AuthPost(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	userData := &struct {
		User *models.User `json:"user"`
	}{User: &models.User{}}

	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		c.config.GetLogger().Error("Error decoding JSON", err)
		return
	}
	svc := c.newService().(*service.AuthService)
	f := svc.GetServiceFuncs()[service.AuthPost]
	if _, err := f.Handler(userData.User); err != nil {
		ErrorHandler(c.config.GetLogger(), w, errorMessage, err, http.StatusBadRequest)
		return
	}

	// Encode the response to JSON
	SetResponseJson(w, userData)
}

func (c *AuthControllerV1) AddHandler(handler ControllerHandlerT) {
	c.httpHandlers = append(c.httpHandlers, ControllerHandlerT{Method: handler.Method, Handler: handler.Handler, Path: handler.Path})
}

func (c *AuthControllerV1) AddHandlers(handlers ...ControllerHandlerT) {
	for _, handler := range handlers {
		c.AddHandler(handler)
	}
}

func (c *AuthControllerV1) GetHandlers() HttpHandlersT {
	return c.httpHandlers
}
