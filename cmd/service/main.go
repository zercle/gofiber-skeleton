package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/samber/do/v2"
	"github.com/zercle/template-go-fiber/internal/config"
	"github.com/zercle/template-go-fiber/internal/handlers"
	"github.com/zercle/template-go-fiber/internal/middleware"
	"github.com/zercle/template-go-fiber/internal/repositories"
	"github.com/zercle/template-go-fiber/internal/usecases"
	_ "github.com/zercle/template-go-fiber/docs"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// Initialize dependency injection
	injector, err := config.InitializeDI(cfg)
	if err != nil {
		log.Fatalf("failed to initialize dependency injection: %v", err)
	}

	// Create structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Go Fiber Microservice Template",
	})

	// Register global middleware
	app.Use(middleware.RecoveryMiddleware(logger))
	app.Use(middleware.CORSMiddleware(cfg))
	app.Use(middleware.LoggerMiddleware(logger))
	app.Use(middleware.RateLimitMiddleware(middleware.NewRateLimiter(100, time.Minute)))

	// Initialize user repository and usecase from DI
	userRepo := do.MustInvoke[*repositories.UserRepository](injector)
	userUsecase := usecases.NewUserUsecase(userRepo)
	userHandler := handlers.NewUserHandler(userUsecase)

	// Health check endpoint (public)
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"message": "Service is running",
		})
	})

	// Swagger documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Public routes (no auth required)
	public := app.Group("/api")
	public.Post("/users/register", userHandler.RegisterUser)

	// Protected routes (auth required)
	protected := app.Group("/api")
	protected.Use(middleware.AuthMiddleware(cfg))
	protected.Get("/users", userHandler.ListUsers)
	protected.Get("/users/email", userHandler.GetUserByEmail)
	protected.Get("/users/:id", userHandler.GetUser)
	protected.Put("/users/:id", userHandler.UpdateUser)
	protected.Delete("/users/:id", userHandler.DeleteUser)

	// Start server in a goroutine
	go func() {
		addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
		log.Printf("Starting server on %s", addr)
		if err := app.Listen(addr); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	<-sigChan

	log.Println("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Printf("error during shutdown: %v", err)
	}

	log.Println("Server stopped")
}
