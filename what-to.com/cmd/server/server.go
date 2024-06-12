package main

import (
	"fmt"

	"what-to.com/pkg/whattoapp"
)

func main() {
	app := whattoapp.NewWhattoApp()
	err := app.Start()
	if err != nil {
		fmt.Printf("Application finished with error: %s", err)
	} else {
		fmt.Println("Application finished successfully")
	}
}
