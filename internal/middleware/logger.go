package middleware

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Logger() fiber.Handler {
	return logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path} - ${ip} ${ua}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
		TimeInterval: 500 * time.Millisecond,
		Output: os.Stdout,
	})
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}