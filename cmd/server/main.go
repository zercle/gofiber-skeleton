package main

import (
	"log"

	"gofiber-skeleton/internal/app"
)

func main() {
	application, err := app.NewApplication()
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}

	application.Start()
	application.WaitForShutdown()
}
