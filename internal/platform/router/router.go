package router

import (
	urlHandler "gofiber-skeleton/internal/url/delivery/http"   // Updated import
	userHandler "gofiber-skeleton/internal/user/delivery/http" // Updated import

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
)

// RegisterRoutes registers the routes for the application.
func RegisterRoutes(app *fiber.App, userHandler *userHandler.UserHandler, urlHandler *urlHandler.URLHandler, jwtSecret string) {
	api := app.Group("/api")

	// User routes
	user := api.Group("/users")
	user.Post("/register", userHandler.Register)
	user.Post("/login", userHandler.Login)

	// URL routes
	url := api.Group("/urls")
	url.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(jwtSecret),
	}))
	url.Post("/", urlHandler.CreateURL)

	// Redirect route
	app.Get("/:shortCode", urlHandler.GetOriginalURL)
}
