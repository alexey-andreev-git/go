package whattoapp

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"what-to.com/internal/config"
	"what-to.com/internal/controller"
	"what-to.com/internal/middleware"
	"what-to.com/internal/migration"
	"what-to.com/internal/models"
	"what-to.com/internal/repository"
	"what-to.com/internal/router"
	"what-to.com/internal/service"
)

type WhatToApp struct {
	appConfig     *config.Config
	appRepository repository.Repository
	appRouter     router.Router
	appMiddleware middleware.Middleware
	httpServer    *http.Server
	httpConfig    config.ConfigT
}

func NewWhattoApp() *WhatToApp {
	app := &WhatToApp{}
	app.appConfig = config.NewConfig()
	// app.appRepository = repository.NewGormPgRepository(app.appConfig)
	// app.appRepository = repository.NewPgRepository(app.appConfig)
	// app.appRouter = router.NewEntityRouter()
	// app.httpConfig = app.appConfig.GetConfig()["http"].(config.ConfigT)
	// app.appMiddleware = middleware.NewEntityMiddleware(app.appConfig, app.appRouter.GetMuxRouter())
	// app.httpServer = &http.Server{
	// 	Addr:    fmt.Sprintf(":%d", app.httpConfig["port"].(int)), // Configure the bind address.
	// 	Handler: app.appMiddleware.GetHandler(),                   // Http handlers here.
	// }
	app.tests()
	return app
}

func (app *WhatToApp) Start() error {
	return nil
	app.appRouter.AddController(
		"entity",
		controller.NewHttpControllerV1(
			app.appConfig,
			service.NewEntityService(app.appConfig, app.appRepository),
		),
	)
	app.appRouter.AddController(
		"auth",
		controller.NewHttpControllerV1(
			app.appConfig,
			service.NewAuthService(app.appConfig, app.appRepository),
		),
	)
	app.appRouter.AddController(
		"front_routes",
		controller.NewHttpControllerV1(
			app.appConfig,
			service.NewFrontRoutesService(app.appConfig, app.appRepository),
		),
	)
	app.appRouter.AddController(
		"front",
		controller.NewHttpControllerV1(
			app.appConfig,
			service.NewFrontService(app.appConfig, app.appRepository),
		),
	)
	// app.tests()
	return app.startServer()
}

func (app *WhatToApp) tests() {
	// u := models.User{
	// 	Name:     "test",
	// 	Password: "test",
	// 	Person: models.Person{
	// 		FirstName:  "first_test",
	// 		MiddleName: "middle_test",
	// 		LastName:   "last_test",
	// 		AddressList: []models.Address{
	// 			{
	// 				Unit:     "unit_test",
	// 				Building: "building_test",
	// 				Street:   "street_test",
	// 				City:     "city_test",
	// 				State:    "state_test",
	// 				ZipCode:  "zip_test",
	// 			},
	// 		},
	// 		EmailsList: []models.Email{
	// 			{
	// 				Name:  "email_test",
	// 				Email: "my@email.net",
	// 			},
	// 		},
	// 		PhonesList: []models.Phone{
	// 			{
	// 				Name:  "phone_test",
	// 				Phone: "1234567890",
	// 			},
	// 		},
	// 	},
	// }
	// app.appRepository.Create(context.Background(), &u)
	// ul := []models.User{}
	// app.appRepository.FindAll(context.Background(), &ul)
	// json, _ := json.Marshal(ul)
	// fmt.Println(json)
	qry := generateCreateTableQuery("addresses", models.Address{})
	fmt.Println(qry)
	qry = generateCreateTableQuery("people", models.Person{})
	fmt.Println(qry)
}

// Generate SQL CREATE TABLE query using structure name
func generateCreateTableQuery(tableName string, model interface{}) string {
	pm := migration.NewPgMigration()
	pm.AddTable(tableName, model)
	// columns := getModelFields(model)
	columns := pm.GetTables()[tableName]
	strColumns := ""
	delimiter := ""
	for _, column := range columns {
		strColumns += delimiter + column.Name + " " + column.Type + " " + strings.Join(column.Tags, " ")
		delimiter = ", "
	}

	// columnsSQL := strings.Join(columns, ", ")
	query := fmt.Sprintf("CREATE TABLE %s IF NOT EXISTS (%s);", tableName, strColumns)
	// query := pm.GetTables()[tableName]

	return query
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
