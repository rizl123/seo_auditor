package main

import (
	"backend/internal/bootstrap"
	"log"
)

func main() {
	app, err := bootstrap.NewApp()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	app.Run()
}
