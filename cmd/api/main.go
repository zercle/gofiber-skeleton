package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	configs "gofiber-skeleton/internal/configs" // Aliased configs
	"gofiber-skeleton/internal/platform/db"
	routerPkg "gofiber-skeleton/internal/platform/router" // Aliased router
	urlHandler "gofiber-skeleton/internal/url/delivery/http"
	urlRepo "gofiber-skeleton/internal/url/repository"
	urlUseCase "gofiber-skeleton/internal/url/usecase"
	userHandler "gofiber-skeleton/internal/user/delivery/http"
	userRepo "gofiber-skeleton/internal/user/repository"
	userUseCase "gofiber-skeleton/internal/user/usecase"

	_ "gofiber-skeleton/api" // Import generated docs

	"github.com/gofiber/fiber/v2"
	swagger "github.com/gofiber/swagger"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9" // Import for Valkey client
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Load configuration
	cfg, err := configs.LoadConfig(os.Getenv("GO_ENV"))
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
		DB:       int(cfg.Cache.DB),
	})

	// Create repositories
	queries := db.New(dbpool)
	userRepository := userRepo.NewSQLUserRepository(queries)
	urlRepository := urlRepo.NewSQLURLRepository(queries, redisClient)

	// Create use cases
	userUsecase := userUseCase.NewUserUseCase(userRepository, cfg.JWT.Secret, cfg.JWT.Expiration)
	urlUsecase := urlUseCase.NewURLUseCase(urlRepository, redisClient)

	// Create handlers
	userHttpHandler := userHandler.NewHTTPUserHandler(userUsecase)
	urlHttpHandler := urlHandler.NewHTTPURLHandler(urlUsecase)

	// Create a new Fiber instance
	app := fiber.New()

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Swagger documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Register routes
	routerPkg.RegisterRoutes(app, userHttpHandler, urlHttpHandler, cfg.JWT.Secret)

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
