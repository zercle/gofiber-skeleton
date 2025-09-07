package middleware

import "github.com/gofiber/fiber/v2"

// Middleware wrapper types to distinguish them in dependency injection
type LoggerMiddleware fiber.Handler
type RecoverMiddleware fiber.Handler
type CORSMiddleware fiber.Handler
type AuthMiddleware fiber.Handler
