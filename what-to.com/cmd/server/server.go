package main

import (
	"fmt"
	"net/http"

	"what-to.com/internal/config"
	"what-to.com/internal/logger"
	"what-to.com/internal/repository"
	"what-to.com/internal/router"
)

func main() {
	appConfig := config.NewConfig("pg_db_connection.yaml")

	appRepository := repository.NewPgRepository(appConfig.GetConfig())
	fmt.Println(appRepository.GetRepoConfigStr())

	appRouter := router.SetupRouter()

	logger.CustomLogger.Fatal("Start server failed:", http.ListenAndServe(":8089", appRouter))
}

// initRepository initializes the repository
// func initApplication() {
// 	// config.ReadConfig()
// 	// repository.SetDBConfig(config.Config["database"].(config.ConfigT))
// 	// repository.ConnectToDB()
// 	appRepo := repository.NewPgRepository()
// }
