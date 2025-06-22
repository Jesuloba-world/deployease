package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/uptrace/bunrouter"

	"github.com/Jesuloba-world/deployease/backend/internal/api"
	"github.com/Jesuloba-world/deployease/backend/internal/app/middleware"
	"github.com/Jesuloba-world/deployease/backend/internal/config"
)

type App struct {
	config *config.Config
	server *http.Server
	router *bunrouter.Router
	api    *api.API
}

func NewApp(cfg *config.Config) *App {
	router := bunrouter.New()

	apiInstance := api.NewAPI(*cfg, router)

	return &App{
		config: cfg,
		router: router,
		api:    apiInstance,
	}
}

func (a *App) Run() error {
	a.setupMiddlewares()

	a.api.InitializeAndRegisterRoutes()

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

func (a *App) setupMiddlewares() {
	recovererConfig := middleware.DefaultRecovererConfig()
	a.router.Use(middleware.Recoverer(recovererConfig))

	timeoutConfig := middleware.DefaultTimeoutConfig()
	a.router.Use(middleware.Timeout(timeoutConfig))

	requestIDCOnfig := middleware.DefaultRequestIDConfig()
	a.router.Use(middleware.RequestID(requestIDCOnfig))

	loggingConfig := middleware.DefaultLoggingConfig()
	a.router.Use(middleware.Logging(loggingConfig))

	corsConfig := middleware.DefaultCORSConfig()
	a.router.Use(middleware.CORS(corsConfig))
}
