package service

import (
	"net/http"

	"github.com/gorilla/mux"
	"what-to.com/internal/config"
	"what-to.com/internal/resources"
)

type (
	EntityService interface {
		NewRestController()
	}
)

// func NewEntityService(appConfig *config.Config) *RestController {
// 	c := &RestController{
// 		httpHandlers: make(HttpHandlersT),
// 	}
// 	c.appRes = resources.NewAppSources()
// 	c.httpHandlers["/entity/{rest:.*}"] = ControllerHandlerT{
// 		Method:  "GET",
// 		Handler: c.entityHandler,
// 	}
// 	c.httpHandlers["/api/{rest:.*}"] = ControllerHandlerT{
// 		Method:  "GET",
// 		Handler: c.apiHandler,
// 	}
// 	return c
// }

func EntityServiceFunction(r *http.Request, config *config.Config) string {
	// Here you would call your repository functions and implement business logic

	muxVars := mux.Vars(r)
	rest := muxVars["rest"]

	appRes := resources.NewAppSources()
	data, err := appRes.GetRes().ReadFile(config.InitDbFileName) // this is the embed.FS
	if err != nil {
		config.GetLogger().Fatal("File read error "+config.InitDbFileName, err)
	}

	// Example: return r *http.Request as a string
	return ("Result: the entity\n" + r.RequestURI + "\n" + rest + "\n" + string(data))

	// return "Result: the entity"
}
