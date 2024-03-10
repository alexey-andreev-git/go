package main

import (
	"log"
	"net/http"

	"what-to.com/internal/router"
)

func main() {
	r := router.SetupRouter()

	log.Fatal(http.ListenAndServe(":8089", r))
}
