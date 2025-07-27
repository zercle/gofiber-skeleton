package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"gofiber-skeleton/internal/configs"
	"gofiber-skeleton/internal/delivery/http"
	"gofiber-skeleton/internal/repository"
	db "gofiber-skeleton/internal/repository/db" // Added import
	"gofiber-skeleton/internal/usecases"

	"github.com/gofiber/fiber/v2"
	swagger "github.com/gofiber/swagger"
	_ "gofiber-skeleton/api" // Import generated docs

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9" // Import for Valkey client
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Load configuration
	cfg, err := configs.LoadConfig()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Connect to the database
	dbpool, err := pgxpool.New(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName, cfg.Database.SSLMode))
	if err != nil {
		slog.Error("Unable to connect to database", "error", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	// Initialize Redis client for Valkey
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Cache.Host, cfg.Cache.Port),
		Password: cfg.Cache.Password,
		DB:       cfg.Cache.DB,
	})

	// Create repositories
	queries := db.New(dbpool)
	userRepo := repository.NewUserRepository(queries)
	urlRepo := repository.NewURLRepository(queries, &repository.RealRedisClient{Client: redisClient})

	// Create use cases
	userUseCase := usecases.NewUserUseCase(userRepo, cfg.JWT.Secret, cfg.JWT.Expiration)
	urlUseCase := usecases.NewURLUseCase(urlRepo, fmt.Sprintf("http://localhost:%d", cfg.Server.Port))

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
			slog.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")

	if err := app.Shutdown(); err != nil {
		slog.Error("Server shutdown failed", "error", err)
		os.Exit(1)
	}

	slog.Info("Server exited gracefully")
}
