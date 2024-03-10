package main

import (
	"net/http"

	"what-to.com/internal/config"
	"what-to.com/internal/logger"
	"what-to.com/internal/repository"
	"what-to.com/internal/router"
)

func main() {
	initApplication()

	r := router.SetupRouter()

	logger.CustomLogger.Fatal("Start server failed:", http.ListenAndServe(":8089", r))
}

// initRepository initializes the repository
func initApplication() {
	config.ReadConfig()
	repository.SetDBConfig(config.Config["database"].(config.ConfigT))
	repository.ConnectToDB()
}
