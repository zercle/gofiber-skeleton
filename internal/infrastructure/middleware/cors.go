package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/zercle/gofiber-skeleton/internal/infrastructure/config"
)

func NewCORS(cfg *config.Config) CORSMiddleware {
	// Convert string to slice if CORS_ORIGINS is a comma-separated string
	origins := cfg.CORS.AllowOrigins
	if len(origins) == 1 && strings.Contains(origins[0], ",") {
		origins = strings.Split(origins[0], ",")
		for i, origin := range origins {
			origins[i] = strings.TrimSpace(origin)
		}
	}

	return CORSMiddleware(cors.New(cors.Config{
		AllowOrigins:     strings.Join(origins, ","),
		AllowMethods:     strings.Join(cfg.CORS.AllowMethods, ","),
		AllowHeaders:     strings.Join(cfg.CORS.AllowHeaders, ","),
		ExposeHeaders:    strings.Join(cfg.CORS.ExposeHeaders, ","),
		AllowCredentials: cfg.CORS.AllowCredentials,
		MaxAge:           cfg.CORS.MaxAge,
	}))
}
