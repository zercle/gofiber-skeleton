package middleware

import (
	"gofiber-skeleton/internal/infra/auth"
	"gofiber-skeleton/pkg/jsend"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Protected creates a new middleware to protect routes using JWT.
func Protected(jwtService auth.JWTService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return jsend.Error(c, "Missing or malformed JWT", http.StatusUnauthorized, http.StatusUnauthorized)
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return jsend.Error(c, "Missing or malformed JWT", http.StatusUnauthorized, http.StatusUnauthorized)
		}

		tokenString := parts[1]
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			return jsend.Error(c, "Invalid or expired JWT", http.StatusUnauthorized, http.StatusUnauthorized)
		}

		// Store user ID in context for downstream handlers
		c.Locals("userID", claims.UserID)

		return c.Next()
	}
}
