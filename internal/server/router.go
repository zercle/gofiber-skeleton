package server

import (
	"fmt"
	"github.com/zercle/gofiber-skeleton/internal/config"
	"github.com/zercle/gofiber-skeleton/internal/logger"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"
)

// SetupRouter initializes the Fiber app and registers routes
func SetupRouter(db *gorm.DB, cfg *config.Config) *fiber.App {
	app := fiber.New()

	app.Use(recover.New())

	// Health check route
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Placeholder route groups
	api := app.Group("/api")

	v1 := api.Group("/v1")
	v1.Get("/auth", stubHandler)     // TODO: Implement authentication routes
	v1.Get("/users", stubHandler)    // TODO: Implement user routes
	v1.Get("/threads", stubHandler)  // TODO: Implement thread routes
	v1.Get("/posts", stubHandler)    // TODO: Implement post routes
	v1.Get("/comments", stubHandler) // TODO: Implement comment routes

	// Global error handler
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(http.StatusNotFound).SendString("Not Found")
	})

	return app
}

// stubHandler is a placeholder for unimplemented routes
func stubHandler(c *fiber.Ctx) error {
	logger.GetLogger().Warn().Msgf("Route %s not implemented", c.OriginalURL())
	return c.Status(fiber.StatusNotImplemented).SendString(fmt.Sprintf("Route %s Not Implemented", c.OriginalURL()))
}