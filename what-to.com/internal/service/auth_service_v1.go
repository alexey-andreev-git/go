package service

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"what-to.com/internal/config"
	"what-to.com/internal/models"
	"what-to.com/internal/repository"
)

// Struct of entity service
type (
	Climes struct {
		Username string `json:"username"`
		jwt.RegisteredClaims
	}
	AuthService struct {
		appRepository repository.Repository
		appConfig     *config.Config
		serviceFuncs  map[RequestType]ServiceFunc
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
		if r.Method == http.MethodGet {
			token := r.Header.Get("Authorization")
			tokenLen := 6
			if strings.ToLower(token[:tokenLen]) != "bearer" {
				ErrorHandler(s.appConfig.GetLogger(), w, errorMessage, fmt.Errorf("incorrect token"), http.StatusBadRequest)
				return
			}
			for tokenLen < len(token) && token[tokenLen] == ' ' {
				tokenLen++
			}
			token = token[tokenLen:] // remove "Token " from token
			bodyJson["user"] = map[string]interface{}{"token": token}
		}
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

func (s *AuthService) V1AuthServiceGet(data interface{}) ([]byte, error) {
	user := data.(*models.User)
	response := ([]byte)(nil)
	if ok, err := ValidateToken(user.Token); !ok {
		return nil, err
	}
	response = ([]byte)("{\"user\": {\"token\": \"" + user.Token + "\"}}")
	s.appConfig.GetLogger().Info("User found:" + string(response))
	return nil, nil
}

func (s *AuthService) V1AuthServicePost(data interface{}) ([]byte, error) {
	user := data.(*models.User)
	token, err := GenerateToken(user.Name)
	if err != nil {
		return nil, err
	}
	if user.Password != "" {
		user.Password = "********"
	}
	user.Token = token
	return nil, nil
}

func (s *AuthService) V1AuthServicePut(bodyJson interface{}) ([]byte, error) {
	result, err := s.appRepository.(*repository.PgRepository).UpdateEntity(bodyJson.(map[string]interface{}))
	if err != nil {
		return nil, err
	}
	rows, rerr := result.RowsAffected()
	return []byte(fmt.Sprintf(jsonOperationResultMsg, "updated", rows, rerr)), nil
}

func (s *AuthService) V1AuthServiceDelete(bodyJson interface{}) ([]byte, error) {
	result, err := s.appRepository.(*repository.PgRepository).DeleteEntity(bodyJson.(map[string]interface{}))
	if err != nil {
		return nil, err
	}
	rows, rerr := result.RowsAffected()
	return []byte(fmt.Sprintf(jsonOperationResultMsg, "deleted", rows, rerr)), nil
}

func (s *AuthService) ValidateUser(user *models.User) error {
	if user.Name == "" {
		return fmt.Errorf("username is empty")
	}
	if user.Password == "" {
		return fmt.Errorf("password is empty")
	}
	return nil
}

func GenerateToken(user string) (string, error) {
	secretKey := []byte("your_secret_key") // Replace with a secure key

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

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (bool, error) {
	secretKey := []byte("your_secret_key") // Replace with a secure key

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return false, err
	}

	// Validate token claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			if time.Unix(int64(exp), 0).Before(time.Now()) {
				return false, fmt.Errorf("token is expired")
			}
		}
		return true, nil
	}

	return false, fmt.Errorf("invalid token")
}
