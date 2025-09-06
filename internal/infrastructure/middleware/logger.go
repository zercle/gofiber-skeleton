package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/config"
)

func NewLogger(cfg *config.Config) fiber.Handler {
	logConfig := logger.Config{
		Next:   nil,
		Done:   nil,
		Format: "[${time}] ${status} - ${method} ${path} - ${latency}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "UTC",
		Output:     os.Stdout,
	}

	if cfg.IsProduction() {
		logConfig.Format = `{"time":"${time}","status":${status},"method":"${method}","path":"${path}","latency":"${latency}","ip":"${ip}","user_agent":"${ua}","error":"${error}"}` + "\n"
	}

	return logger.New(logConfig)
}