package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"


	"gofiber-skeleton/internal/configs"
	"gofiber-skeleton/internal/delivery/http"
	"gofiber-skeleton/internal/repository"
	"gofiber-skeleton/internal/usecases"

	"github.com/gofiber/fiber/v2"
	swagger "github.com/gofiber/swagger"
	_ "gofiber-skeleton/api" // Import generated docs

	"github.com/jackc/pgx/v5/pgxpool"
)

// @title Go Fiber URL Shortener API
// @version 1.0
// @description This is a sample URL shortener service.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
func main() {
	// Load configuration
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to the database
	dbpool, err := pgxpool.New(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName, cfg.Database.SSLMode))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()

	// Create repositories
	userRepo := repository.NewUserRepository(dbpool)
	urlRepo := repository.NewURLRepository(dbpool)

	// Create use cases
	userUseCase := usecases.NewUserUseCase(userRepo, cfg.JWT.Secret, cfg.JWT.Expiration)
	urlUseCase := usecases.NewURLUseCase(urlRepo)

	// Create handlers
	userHandler := http.NewUserHandler(userUseCase)
	urlHandler := http.NewURLHandler(urlUseCase)

	// Create a new Fiber instance
	app := fiber.New()

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Swagger documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Register routes
	http.RegisterRoutes(app, userHandler, urlHandler, cfg.JWT.Secret)

	// Start the server in a goroutine
	go func() {
		if err := app.Listen(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	if err := app.Shutdown(); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server exited gracefully")
}
