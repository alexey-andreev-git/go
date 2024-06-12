package whattoapp

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"what-to.com/internal/config"
	"what-to.com/internal/controller"
	"what-to.com/internal/repository"
	"what-to.com/internal/router"
	"what-to.com/internal/service"
)

type WhatToApp struct {
	appConfig     *config.Config
	appRepository repository.Repository
	appRouter     router.Router
	appService    service.Service
	httpServer    *http.Server
	httpConfig    config.ConfigT
}

func NewWhattoApp() *WhatToApp {
	app := &WhatToApp{}
	app.appConfig = config.NewConfig()
	app.appRepository = repository.NewPgRepository(app.appConfig)
	app.appRouter = router.NewEntityRouter()
	app.appService = service.NewEntityService(app.appConfig, app.appRepository)
	app.httpConfig = app.appConfig.GetConfig()["http"].(config.ConfigT)
	app.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", app.httpConfig["port"].(int)), // Configure the bind address.
		Handler: app.appRouter.GetMuxRouter(),                     // Http handlers here.
	}
	return app
}

func (app *WhatToApp) Start() error {
	app.appRouter.AddController("front", controller.NewFrontController(app.appConfig))
	app.appRouter.AddController("rest", controller.NewRestController(app.appConfig, app.appService))
	return app.startServer()
}

func (app *WhatToApp) startServer() (err error) {
	// Create a context for graceful shutdown.
	ctxNotifySignal, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start the server in a separate goroutine.
	go func() {
		app.appConfig.GetLogger().Info(fmt.Sprintf("Starting server on : <%v>...", app.httpServer.Addr))
		err = app.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			app.appConfig.GetLogger().Fatal("Start server failed:", err)
		}
	}()

	// Wait for the context to be canceled (e.g., by a signal).
	<-ctxNotifySignal.Done()

	app.appConfig.GetLogger().Info("Shutdown signal received")

	// Create a new context with a timeout for the graceful shutdown.
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server.
	err = app.httpServer.Shutdown(ctxShutdown)
	if err != nil {
		app.appConfig.GetLogger().Fatal("Server forced to shutdown:", err)
	} else {
		app.appConfig.GetLogger().Info("Server shutdown gracefully")
	}
	return err
}
