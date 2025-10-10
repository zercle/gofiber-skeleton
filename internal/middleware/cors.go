package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// CORS returns CORS middleware with default configuration
func CORS() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: false,
		ExposeHeaders:    "Content-Length",
		MaxAge:           86400,
	})
}

// CORSWithConfig returns CORS middleware with custom configuration
func CORSWithConfig(allowOrigins, allowMethods, allowHeaders string, allowCredentials bool) fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     allowMethods,
		AllowHeaders:     allowHeaders,
		AllowCredentials: allowCredentials,
		ExposeHeaders:    "Content-Length",
		MaxAge:           86400,
	})
}
