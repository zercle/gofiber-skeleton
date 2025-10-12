package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/zercle/gofiber-skeleton/internal/config"
	"github.com/zercle/gofiber-skeleton/internal/handlers"
	"github.com/zercle/gofiber-skeleton/internal/middleware"
	"github.com/zercle/gofiber-skeleton/pkg/logger"
)

// Server represents the HTTP server
type Server struct {
	app        *fiber.App
	config     *config.Config
	healthHandler *handlers.HealthHandler
	apiHandler   *handlers.APIHandler
}

// NewServer creates a new server instance
func NewServer(cfg *config.Config) *Server {
	app := fiber.New(fiber.Config{
		AppName:      "Go Fiber Backend",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	})

	// Create handlers
	version := "1.0.0" // You can get this from build info
	healthHandler := handlers.NewHealthHandler(version)
	apiHandler := handlers.NewAPIHandler()

	server := &Server{
		app:          app,
		config:       cfg,
		healthHandler: healthHandler,
		apiHandler:   apiHandler,
	}

	// Setup routes and middleware
	server.setupRoutes()
	server.setupMiddleware()

	return server
}

// setupMiddleware configures all middleware
func (s *Server) setupMiddleware() {
	middleware.SetupMiddleware(s.app, s.config)
}

// setupRoutes configures all application routes
func (s *Server) setupRoutes() {
	// Health check routes
	health := s.app.Group("/health")
	health.Get("/", s.healthHandler.Health)
	health.Get("/live", s.healthHandler.Liveness)
	health.Get("/ready", s.healthHandler.Readiness)

	// API routes
	api := s.app.Group("/api/v1")

	// Todo routes
	todos := api.Group("/todos")
	todos.Get("/", s.apiHandler.GetTodos)
	todos.Get("/:id", s.apiHandler.GetTodo)
	todos.Post("/", s.apiHandler.CreateTodo)
	todos.Put("/:id", s.apiHandler.UpdateTodo)
	todos.Delete("/:id", s.apiHandler.DeleteTodo)

	// Stats endpoint
	api.Get("/stats", s.apiHandler.GetStats)

	// Root endpoint
	s.app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(map[string]interface{}{
			"message":   "Welcome to Go Fiber Backend API",
			"version":   "1.0.0",
			"timestamp": time.Now(),
			"docs":      "/health",
		})
	})
}

// Start starts the server with graceful shutdown
func (s *Server) Start() error {
	// Create a channel to receive OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		addr := s.config.GetAddress()
		logger.Info("Starting server",
			"address", addr,
			"environment", s.config.Server.Environment,
		)

		if err := s.app.Listen(addr); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to start server", "error", err)
			quit <- syscall.SIGTERM
		}
	}()

	// Wait for interrupt signal
	<-quit
	logger.Info("Shutting down server...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), s.config.Server.ShutdownTimeout)
	defer cancel()

	// Attempt graceful shutdown
	if err := s.app.ShutdownWithContext(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
		return err
	}

	logger.Info("Server shutdown completed")
	return nil
}

// GetApp returns the fiber app instance (for testing purposes)
func (s *Server) GetApp() *fiber.App {
	return s.app
}

// GetConfig returns the server configuration
func (s *Server) GetConfig() *config.Config {
	return s.config
}

// HealthCheck performs a quick health check on the server
func (s *Server) HealthCheck() error {
	// Create a simple HTTP client
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Make a request to the health endpoint
	url := fmt.Sprintf("http://%s/health/live", s.config.GetAddress())
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check returned status: %d", resp.StatusCode)
	}

	return nil
}