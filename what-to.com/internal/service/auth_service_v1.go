package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"what-to.com/internal/config"
	"what-to.com/internal/repository"
)

// Struct of entity service
type (
	AuthService struct {
		appRepository repository.Repository
		appConfig     *config.Config
		serviceFuncs  map[RequestType]ServiceFunc
	}
	User struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Token    string `json:"token"`
	}
	Climes struct {
		Username string `json:"username"`
		jwt.RegisteredClaims
	}
)

const (
	authPath = "/auth"
)

const (
	AuthGet RequestType = iota
	AuthPost
	AuthPut
	AuthDelete
)

func NewAuthService(appConfig *config.Config, appRepo repository.Repository) *AuthService {
	s := &AuthService{
		appConfig:     appConfig,
		appRepository: appRepo,
	}
	s.registerServiceFuncs()
	return s
}

func (s *AuthService) registerServiceFuncs() {
	s.serviceFuncs = map[RequestType]ServiceFunc{
		AuthGet:    {s.V1AuthServiceGet, "GET", apiV1Path + authPath + restWildcardPath},
		AuthPost:   {s.V1AuthServicePost, "POST", apiV1Path + authPath + restWildcardPath},
		AuthPut:    {s.V1AuthServicePut, "PUT", apiV1Path + authPath + restWildcardPath},
		AuthDelete: {s.V1AuthServiceDelete, "DELETE", apiV1Path + authPath + restWildcardPath},
	}
}

func (s *AuthService) GetServiceFuncs() map[RequestType]ServiceFunc {
	return s.serviceFuncs
}

func (s *AuthService) ServiceFunction(w http.ResponseWriter, r *http.Request, ver string, reqType RequestType) {
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
		// ErrorHandler(s.appConfig.GetLogger(), w, errorMessage, err, http.StatusBadRequest)
		// return
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

func (s *AuthService) V1AuthServiceGet(bodyJson map[string]interface{}) ([]byte, error) {
	if bodyJson["email"] != "test@email.net" {
		return nil, fmt.Errorf("user not found")
	}
	response := ([]byte)("{\"message\":\"User found\"}")
	return response, nil
}

func (s *AuthService) V1AuthServicePost(bodyJson map[string]interface{}) ([]byte, error) {
	// response := ([]byte)("{\"message\":\"User created\"}")
	user := bodyJson["user"].(map[string]interface{})
	userName := user["username"].(string)
	token, err := GenerateToken(userName)
	if err != nil {
		return nil, err
	}
	response := ([]byte)("{\"user\": {\"token\": \"" + token + "\"}}")
	s.appConfig.GetLogger().Info("User created:" + string(response))
	return response, nil
}

func (s *AuthService) V1AuthServicePut(bodyJson map[string]interface{}) ([]byte, error) {
	result, err := s.appRepository.(*repository.PgRepository).UpdateEntity(bodyJson)
	if err != nil {
		return nil, err
	}
	rows, rerr := result.RowsAffected()
	return []byte(fmt.Sprintf(jsonOperationResultMsg, "updated", rows, rerr)), nil
}

func (s *AuthService) V1AuthServiceDelete(bodyJson map[string]interface{}) ([]byte, error) {
	result, err := s.appRepository.(*repository.PgRepository).DeleteEntity(bodyJson)
	if err != nil {
		return nil, err
	}
	rows, rerr := result.RowsAffected()
	return []byte(fmt.Sprintf(jsonOperationResultMsg, "deleted", rows, rerr)), nil
}

func GenerateToken(user string) (string, error) {
	jwtKey := []byte("your_secret_key") // Replace with a secure key

	period := 24

	expirationTime := time.Now().Add(time.Duration(period) * time.Hour)

	claims := &Climes{
		Username: user,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:  "what-to.com",
			Subject: "authorization",
			// Audience:  jwt.Audience{"what-to.com"},
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
