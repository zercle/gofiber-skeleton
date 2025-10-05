package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// CORS creates a CORS middleware with production-ready settings
func CORS() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://localhost:8080",
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length,Content-Range",
		MaxAge:           3600,
	})
}
