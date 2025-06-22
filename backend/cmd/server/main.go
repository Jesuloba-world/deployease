package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/Jesuloba-world/deployease/backend/internal/app"
	"github.com/Jesuloba-world/deployease/backend/internal/config"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or could not be loaded: %v", err)
	}

	// load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// load application
	application := app.NewApp(cfg)

	if err := application.Run(); err != nil {
		log.Fatalf("Application failed to start: %v", err)
	}
}
