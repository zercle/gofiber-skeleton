package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zercle/gofiber-skelton/internal/handlers"
)

// SetupRoutes is the Router for GoFiber App
func (r *RouterResources) SetupRoutes(app *fiber.App) {

	app.Get("/", handlers.Index())

	groupApiV1 := app.Group("/api/v:version?", apiLimiter)
	{
		groupApiV1.Get("/", handlers.Index())
	}
}
