package main

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"what-to.com/internal/config"
	"what-to.com/internal/controller"
	"what-to.com/internal/logger"
	"what-to.com/internal/repository"
	"what-to.com/internal/router"
	"what-to.com/internal/service"
)

func main() {
	appConfig := config.NewConfig()
	appRepository := repository.NewPgRepository(appConfig)
	fmt.Println(appRepository.GetRepoConfigStr())
	appRouter := router.NewEntityRouter()
	appService := service.NewEntityService(appConfig, appRepository)
	appRouter.AddController("front", controller.NewFrontController(appConfig))
	appRouter.AddController("rest", controller.NewRestController(appConfig, appService))

	httpConfig := appConfig.GetConfig()["http"].(config.ConfigT)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", httpConfig["port"].(int)), // Configure the bind address.
		Handler: appRouter.GetMuxRouter(),                     // Http handlers here.
	}
	startServer(server, appConfig.GetLogger())
}

func startServer(server *http.Server, clogger logger.Logger) {
	// Create a context for graceful shutdown.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start the server in a separate goroutine.
	go func() {
		clogger.Info(fmt.Sprintf("Starting server on : <%v>...", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			clogger.Fatal("Start server failed:", err)
		}
	}()

	// Wait for the context to be canceled (e.g., by a signal).
	<-ctx.Done()
	clogger.Info("Shutdown signal received")

	// Create a new context with a timeout for the graceful shutdown.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server.
	if err := server.Shutdown(shutdownCtx); err != nil {
		clogger.Fatal("Server forced to shutdown:", err)
	} else {
		clogger.Info("Server shutdown gracefully")
	}
}
