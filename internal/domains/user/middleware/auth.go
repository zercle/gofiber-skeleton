package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/usecase"
	"github.com/zercle/gofiber-skeleton/internal/response"
)

// JWTAuth returns a JWT authentication middleware
func JWTAuth(authUsecase usecase.AuthUsecase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Unauthorized(c, "Missing authorization header")
		}

		// Check if it's a Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return response.Unauthorized(c, "Invalid authorization format. Use: Bearer <token>")
		}

		token := parts[1]
		if token == "" {
			return response.Unauthorized(c, "Missing token")
		}

		// Verify token
		userID, err := authUsecase.VerifyToken(token)
		if err != nil {
			return response.Unauthorized(c, "Invalid or expired token")
		}

		// Store user ID in context
		c.Locals("userID", userID.String())

		return c.Next()
	}
}

// OptionalJWTAuth returns an optional JWT authentication middleware
// It tries to authenticate but doesn't fail if no token is provided
func OptionalJWTAuth(authUsecase usecase.AuthUsecase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		// Check if it's a Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Next()
		}

		token := parts[1]
		if token == "" {
			return c.Next()
		}

		// Verify token
		userID, err := authUsecase.VerifyToken(token)
		if err != nil {
			return c.Next()
		}

		// Store user ID in context
		c.Locals("userID", userID.String())

		return c.Next()
	}
}
