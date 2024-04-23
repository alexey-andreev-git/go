package service

import (
	"net/http"

	"github.com/gorilla/mux"
	"what-to.com/internal/config"
	"what-to.com/internal/resources"
)

func EntityServiceFunction(r *http.Request, c *config.Config) string {
	// Here you would call your repository functions and implement business logic

	muxVars := mux.Vars(r)
	rest := muxVars["rest"]

	appRes := resources.NewAppSources()
	// data, err := appRes.GetRes().ReadFile(c.GetConfig()["initDbFileName"]) // this is the embed.FS
	data, err := appRes.GetRes().ReadFile(c.GetConfig()[config.KeyInitDbFileName].(string)) // this is the embed.FS
	if err != nil {
		// c.GetLogger().Fatal(fmt.Sprintf("File read error [%s] [%s]", c.GetConfig()["initDbFileName"], err))
		c.GetLogger().Fatal("File read error [%s] "+c.GetConfig()[config.KeyInitDbFileName].(string), err)
	}

	// Example: return r *http.Request as a string
	return ("Result: the entity\n" + r.RequestURI + "\n" + rest + "\n" + string(data))

	// return "Result: the entity"
}
