package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/zercle/gofiber-skeleton/internal/logger"
)

func AuthRequired(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Authorization header missing"})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid Authorization header format"})
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			logger.GetLogger().Error().Err(err).Msg("Invalid JWT token")
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid or expired JWT token"})
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if exp, ok := claims["exp"].(float64); ok {
				if int64(exp) < time.Now().Unix() {
					return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "JWT token is expired"})
				}
			} else {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "JWT token expiration claim missing"})
			}

			userID, err := uuid.Parse(claims["user_id"].(string))
			if err != nil {
				logger.GetLogger().Error().Err(err).Msg("Failed to parse user ID from JWT claims")
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid user ID in token"})
			}
			c.Locals("user_id", userID)
			return c.Next()
		}

		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid JWT token claims"})
	}
}