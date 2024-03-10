package main

import (
	"net/http"

	"what-to.com/internal/logger"
	"what-to.com/internal/repository"
	"what-to.com/internal/router"
)

func main() {
	initRepository()

	r := router.SetupRouter()

	logger.Log.Fatalf("Start server failed: %v", http.ListenAndServe(":8089", r))
}

// initRepository initializes the repository
func initRepository() {
	repository.SetDBConfig(repository.ReadDBConfig())
	repository.ConnectToDB()
}
