package server

import (
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"

	"github.com/zercle/gofiber-skeleton/internal/infrastructure/config"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/middleware"
	"github.com/zercle/gofiber-skeleton/internal/shared/jsend"
)

// NewFiberApp creates a new Fiber application with global middlewares.
func NewFiberApp(
	cfg *config.Config,
	logger middleware.LoggerMiddleware,
	recover middleware.RecoverMiddleware,
	cors middleware.CORSMiddleware,
) *fiber.App {
	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		ErrorHandler: ErrorHandler,
	})

	// Global middlewares
	app.Use(fiber.Handler(logger))
	app.Use(fiber.Handler(recover))
	app.Use(fiber.Handler(cors))

	return app
}

// StartServer hooks the Fiber app to the fx lifecycle for graceful startup and shutdown.
func StartServer(
	lc fx.Lifecycle,
	app *fiber.App,
	cfg *config.Config,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				addr := fmt.Sprintf(":%s", cfg.App.Port)
				log.Printf("Server starting on port %s", cfg.App.Port)
				log.Printf("Environment: %s", cfg.App.Environment)
				log.Printf("Swagger UI: http://localhost:%s/swagger/", cfg.App.Port)

				if err := app.Listen(addr); err != nil {
					log.Fatalf("Failed to start server: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Shutting down server...")
			return app.Shutdown()
		},
	})
}

// ErrorHandler customizes Fiber's error handling to return JSend-compliant responses.
func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	log.Printf("Error: %v", err)

	return jsend.SendFail(c, code, map[string]string{
		"message": message,
	})
}
