package http

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"gofiber-skeleton/internal/auth"
	"gofiber-skeleton/internal/configs"
)

func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "data": "Missing or malformed JWT"})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	cfg, err := configs.LoadConfig()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to load config"})
	}

	claims, err := auth.ValidateToken(tokenString, cfg.JWT.Secret)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "data": "Invalid or expired JWT"})
	}

	c.Locals("user_id", claims.UserID)
	return c.Next()
}
