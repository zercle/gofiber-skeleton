package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/swagger"
	_ "github.com/zercle/gofiber-skeleton/docs"

	"github.com/samber/do"

	"github.com/zercle/gofiber-skeleton/internal/infrastructure/config"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/database"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/di"
	"github.com/zercle/gofiber-skeleton/internal/shared/middleware"
	sharedRouter "github.com/zercle/gofiber-skeleton/internal/shared/router"
	"github.com/zercle/gofiber-skeleton/internal/shared/response"
)

// @title Go Fiber Skeleton API
// @version 1.0
// @description Production-ready Go Fiber backend template with Clean Architecture
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	ctx := context.Background()
	injector := di.NewServiceContainer(ctx)
	defer func() {
		if err := do.Shutdown[*do.Injector](injector); err != nil {
			log.Printf("Failed to shutdown dependency injector: %v", err)
		}
	}()

	cfg := do.MustInvoke[*config.Config](injector)
	do.MustInvoke[*database.Database](injector) // Ensure database is initialized and closed by do.Shutdown
	
	// Handlers are retrieved by the routers themselves via the injector

	app := fiber.New(fiber.Config{
		Prefork:               false,
		CaseSensitive:         true,
		StrictRouting:         true,
		DisableStartupMessage: false,
		ReadTimeout:           time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout:          time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:           time.Duration(cfg.Server.IdleTimeout) * time.Second,
	})

	setupMiddleware(app)
	setupRoutes(app, injector, cfg)

	go func() {
		addr := cfg.GetServerAddr()
		log.Printf("Starting server on %s", addr)
		if err := app.Listen(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	// Database closing is handled by do.Shutdown(injector)

	log.Println("Server exited")
}

func setupMiddleware(app *fiber.App) {
	app.Use(middleware.NewRecovery())
	app.Use(middleware.NewLogger())
	app.Use(middleware.NewRequestID())
	app.Use(middleware.NewCORS())
	app.Use(middleware.NewSecurityHeaders())
	app.Use(middleware.NewRateLimiting())
	app.Use(middleware.NewTimeout(30 * time.Second))

	if os.Getenv("APP_ENV") == "development" {
		app.Use(pprof.New())
	}
}

func setupRoutes(app *fiber.App, injector *do.Injector, cfg *config.Config) {
	app.Get("/health", middleware.NewHealthCheck())
	app.Get("/swagger/*", swagger.New(swagger.Config{}))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to Go Fiber Skeleton API",
			"version": "1.0.0",
			"docs":    "/swagger/index.html",
		})
	})

	api := app.Group("/api/v1")

	routers := do.MustInvoke[[]sharedRouter.Router](injector)
	for _, r := range routers {
		r.RegisterRoutes(api)
	}

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(response.Response{
			Success: false,
			Message: "Route not found",
			Error:   "The requested route does not exist",
		})
	})
}