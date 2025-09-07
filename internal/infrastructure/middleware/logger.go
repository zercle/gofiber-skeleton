package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/zercle/gofiber-skeleton/internal/infrastructure/config"
)

func NewLogger(cfg *config.Config) LoggerMiddleware {
	logFormat := "[${time}] ${status} - ${method} ${path} - ${ip} - ${latency} - ${ua}\n"

	if cfg.IsProduction() {
		// Production format (JSON-like structured logging)
		logFormat = `{"time":"${time}","level":"info","status":${status},"method":"${method}","path":"${path}","ip":"${ip}","latency":"${latency}","user_agent":"${ua}"}` + "\n"
	}

	return LoggerMiddleware(logger.New(logger.Config{
		Format:     logFormat,
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "UTC",
		Output:     os.Stdout,
	}))
}
