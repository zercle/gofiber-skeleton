package http

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
)

// RegisterRoutes registers the routes for the application.
func RegisterRoutes(app *fiber.App, userHandler *UserHandler, urlHandler *URLHandler, jwtSecret string) {
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
	url.Post("/", urlHandler.CreateShortURL)

	// Redirect route
	app.Get("/:shortCode", urlHandler.Redirect)

	// QR code route
	app.Get("/:shortCode/qr", urlHandler.GetQRCode)
}
