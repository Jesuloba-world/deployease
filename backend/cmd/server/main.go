package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/reqlog"

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
	app := NewApp(cfg)

	if err := app.Run(); err != nil {
		log.Fatalf("Application failed to start: %v", err)
	}
}

type App struct {
	config *config.Config
	server *http.Server
	router *bunrouter.Router
}

func NewApp(cfg *config.Config) *App {
	return &App{
		config: cfg,
		router: bunrouter.New(
			bunrouter.Use(reqlog.NewMiddleware()),
		),
	}
}

func (a *App) Run() error {
	a.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%s", a.config.Server.Host, a.config.Server.Port),
		Handler:      a.router,
		ReadTimeout:  a.config.Server.ReadTimeout,
		WriteTimeout: a.config.Server.WriteTimeout,
		IdleTimeout:  a.config.Server.IdleTimeout,
	}

	go func() {
		log.Printf("Starting DeployEase server on %s:%s", a.config.Server.Host, a.config.Server.Port)
		log.Printf("Environment: %s", a.config.Environment)
		log.Printf("Server address: %s", a.server.Addr)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
			log.Fatalf("Failed to start server: %v", err)
		}
		log.Println("Server stopped")
	}()

	return a.gracefulShutdown()
}

func (a *App) gracefulShutdown() error {
	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := a.server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
		return err
	}

	log.Println("Server exited gracefully")
	return nil
}
