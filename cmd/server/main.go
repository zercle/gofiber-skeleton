package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"go.uber.org/fx"

	"github.com/zercle/gofiber-skeleton/internal/domains/auth/handlers"
	authRepos "github.com/zercle/gofiber-skeleton/internal/domains/auth/repositories"
	authRoutes "github.com/zercle/gofiber-skeleton/internal/domains/auth/routes"
	authUsecases "github.com/zercle/gofiber-skeleton/internal/domains/auth/usecases"
	postHandlers "github.com/zercle/gofiber-skeleton/internal/domains/posts/handlers"
	postRepos "github.com/zercle/gofiber-skeleton/internal/domains/posts/repositories"
	postRoutes "github.com/zercle/gofiber-skeleton/internal/domains/posts/routes"
	postUsecases "github.com/zercle/gofiber-skeleton/internal/domains/posts/usecases"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/config"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/database"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/middleware"
	"github.com/zercle/gofiber-skeleton/pkg/utils"
)

// @title Go Fiber Backend API
// @version 1.0
// @description A production-ready Go Fiber backend template with Clean Architecture
// @contact.name API Support
// @contact.email support@example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:3000
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	app := fx.New(
		fx.Provide(
			config.Load,
			database.NewDatabase,
			utils.NewJWTManager,
			authRepos.NewPostgresUserRepository,
			postRepos.NewPostgresPostRepository,
			authUsecases.NewAuthUsecase,
			postUsecases.NewPostUsecase,
			handlers.NewAuthHandler,
			postHandlers.NewPostHandler,
			NewFiberApp,
		),
		fx.Invoke(StartServer),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := app.Start(ctx); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	if err := app.Stop(ctx); err != nil {
		log.Fatalf("Failed to stop application: %v", err)
	}

	log.Println("Server stopped")
}

func NewFiberApp(cfg *config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		ErrorHandler: middleware.ErrorHandler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	})

	app.Use(middleware.NewLogger(cfg))
	app.Use(middleware.NewCORS(cfg))
	app.Use(middleware.NewRecover(cfg))

	if cfg.IsDevelopment() {
		app.Get("/docs/*", swagger.HandlerDefault)
		app.Get("/docs", func(c *fiber.Ctx) error {
			return c.Redirect("/docs/")
		})
	}

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
			"service":   cfg.App.Name,
			"version":   cfg.App.Version,
		})
	})

	return app
}

func StartServer(
	app *fiber.App,
	cfg *config.Config,
	db *database.Database,
	authHandler *handlers.AuthHandler,
	postHandler *postHandlers.PostHandler,
) {
	authRoutes.SetupAuthRoutes(app, authHandler, cfg)
	postRoutes.SetupPostRoutes(app, postHandler, cfg)

	log.Printf("Starting server on port %s in %s mode", cfg.App.Port, cfg.App.Environment)

	if err := app.Listen(":" + cfg.App.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}