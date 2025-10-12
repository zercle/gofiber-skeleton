package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// CORSConfig provides CORS middleware configuration
type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	ExposeHeaders    []string
	MaxAge           int
}

// DefaultCORSConfig returns default CORS configuration
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           86400, // 24 hours
	}
}

// ProductionCORSConfig returns production-ready CORS configuration
func ProductionCORSConfig(allowedOrigins []string) CORSConfig {
	return CORSConfig{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           86400, // 24 hours
	}
}

// CORS creates CORS middleware with custom configuration
func CORS(config ...CORSConfig) fiber.Handler {
	var cfg CORSConfig
	if len(config) > 0 {
		cfg = config[0]
	} else {
		cfg = DefaultCORSConfig()
	}

	return func(c *fiber.Ctx) error {
		origin := c.Get("Origin")

		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range cfg.AllowOrigins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				allowed = true
				break
			}
		}

		if allowed {
			// Set CORS headers
			if len(cfg.AllowOrigins) > 0 {
				if cfg.AllowOrigins[0] == "*" {
					c.Set("Access-Control-Allow-Origin", "*")
				} else {
					c.Set("Access-Control-Allow-Origin", origin)
				}
			}

			if len(cfg.AllowMethods) > 0 {
				c.Set("Access-Control-Allow-Methods", strings.Join(cfg.AllowMethods, ", "))
			}

			if len(cfg.AllowHeaders) > 0 {
				c.Set("Access-Control-Allow-Headers", strings.Join(cfg.AllowHeaders, ", "))
			}

			if len(cfg.ExposeHeaders) > 0 {
				c.Set("Access-Control-Expose-Headers", strings.Join(cfg.ExposeHeaders, ", "))
			}

			if cfg.AllowCredentials {
				c.Set("Access-Control-Allow-Credentials", "true")
			}

			c.Set("Access-Control-Max-Age", fmt.Sprintf("%d", cfg.MaxAge))
		}

		// Handle preflight requests
		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusNoContent)
		}

		return c.Next()
	}
}

// DevelopmentCORS creates permissive CORS for development
func DevelopmentCORS() fiber.Handler {
	return CORS(DefaultCORSConfig())
}

// RestrictedCORS creates restrictive CORS for production
func RestrictedCORS(allowedOrigins []string) fiber.Handler {
	return CORS(ProductionCORSConfig(allowedOrigins))
}