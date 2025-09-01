package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/samber/do/v2"

	_ "github.com/zercle/gofiber-skeleton/docs"
	"github.com/zercle/gofiber-skeleton/internal/config"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/container"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/database"
	"github.com/zercle/gofiber-skeleton/internal/middleware"
	"github.com/zercle/gofiber-skeleton/internal/user/handler"
	authHandler "github.com/zercle/gofiber-skeleton/internal/auth/handler"
)

func main() {
	// Setup dependency injection container
	diContainer, err := container.SetupContainer()
	if err != nil {
		log.Fatalf("Failed to setup DI container: %v", err)
	}
	defer diContainer.Shutdown()

	// Get dependencies from container
	cfg := do.MustInvoke[*config.Config](diContainer)
	db := do.MustInvoke[*database.Database](diContainer)
	userHandler := do.MustInvoke[*handler.UserHandler](diContainer)

	// Create JWT middleware
	jwtMiddleware := middleware.NewJWTMiddleware(cfg)

	app := fiber.New(fiber.Config{
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
		ReadTimeout:  cfg.App.ReadTimeout,
		ErrorHandler: middleware.ErrorHandler(),
		Prefork:      runtime.NumCPU() > 1,
	})

	app.Use(recover.New())
	app.Use(middleware.Logger())
	app.Use(middleware.CORS())

	// Swagger documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	userHandler.InitUserRoutes(app, jwtMiddleware)
	// Init auth routes
	authHandler := do.MustInvoke[*authHandler.AuthHandler](diContainer)
	authHandler.InitAuthRoutes(app)

	// Health check cache
	type healthStatus struct {
		isHealthy  bool
		lastChecked time.Time
	}
	var healthCache struct {
		sync.RWMutex
		status healthStatus
	}
	const cacheDuration = 5 * time.Second

	app.Get("/health", func(c *fiber.Ctx) error {
		healthCache.RLock()
		if time.Since(healthCache.status.lastChecked) < cacheDuration {
			defer healthCache.RUnlock()
			if healthCache.status.isHealthy {
				return c.JSON(fiber.Map{"status": "ok", "message": "Server is healthy (cached)"})
			}
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"status": "error", "message": "Database connection failed (cached)"})
		}
		healthCache.RUnlock()

		healthCache.Lock()
		defer healthCache.Unlock()
		// Double-check in case another goroutine updated the cache while we were waiting for the lock
		if time.Since(healthCache.status.lastChecked) < cacheDuration {
			if healthCache.status.isHealthy {
				return c.JSON(fiber.Map{"status": "ok", "message": "Server is healthy (cached)"})
			}
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"status": "error", "message": "Database connection failed (cached)"})
		}

		err := db.Health()
		healthCache.status.lastChecked = time.Now()
		if err != nil {
			healthCache.status.isHealthy = false
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status":  "error",
				"message": "Database connection failed",
			})
		}

		healthCache.status.isHealthy = true
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Server is healthy",
		})
	})

	go func() {
		log.Printf("🚀 Server starting on port %d", cfg.App.Port)
		log.Printf("📚 Swagger docs available at http://localhost:%d/swagger/", cfg.App.Port)
		if err := app.Listen(fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port)); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Server shutting down...")

	if err := app.Shutdown(); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server stopped")
}
