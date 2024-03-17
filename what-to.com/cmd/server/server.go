package main

import (
	"fmt"
	"net/http"

	"what-to.com/internal/config"
	"what-to.com/internal/controller"
	"what-to.com/internal/logger"
	"what-to.com/internal/repository"
	"what-to.com/internal/router"
)

func main() {
	appConfig := config.NewConfig("pg_db_connection.yaml")

	appRepository := repository.NewPgRepository(appConfig.GetConfig())
	fmt.Println(appRepository.GetRepoConfigStr())

	appRouter := router.NewEntityRouter()
	appRouter.AddController("front", controller.NewFrontController())
	appRouter.AddController("rest", controller.NewRestController())

	logger.CustomLogger.Fatal("Start server failed:", http.ListenAndServe(":8088", appRouter.GetMuxRouter()))
}
