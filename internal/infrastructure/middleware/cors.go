package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/config"
)

func NewCORS(cfg *config.Config) fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     joinStrings(cfg.CORS.AllowOrigins, ","),
		AllowMethods:     joinStrings(cfg.CORS.AllowMethods, ","),
		AllowHeaders:     joinStrings(cfg.CORS.AllowHeaders, ","),
		ExposeHeaders:    joinStrings(cfg.CORS.ExposeHeaders, ","),
		AllowCredentials: cfg.CORS.AllowCredentials,
		MaxAge:           cfg.CORS.MaxAge,
	})
}

func joinStrings(slice []string, separator string) string {
	if len(slice) == 0 {
		return ""
	}
	result := slice[0]
	for i := 1; i < len(slice); i++ {
		result += separator + slice[i]
	}
	return result
}