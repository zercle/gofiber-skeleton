package handler

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"

	"github.com/zercle/gofiber-skeleton/internal/infrastructure/config"
)

// SetupRoutes sets up the routes for the demo handler.
func SetupRoutes(app *fiber.App, cfg *config.Config, demoHandler *DemoHandler) {
	jwtMiddleware := jwtware.New(jwtware.Config{
		SigningKey: []byte(cfg.JWT.Secret),
	})

	demoGroup := app.Group("/api/v1/demo")
	demoGroup.Post("/transaction", jwtMiddleware, demoHandler.PerformTransactionDemo)
	demoGroup.Get("/joined", jwtMiddleware, demoHandler.GetJoinedDataDemo)
}
