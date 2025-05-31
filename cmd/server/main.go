package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gofiber-skeleton/config"
	"gofiber-skeleton/internal/handler"
	"gofiber-skeleton/internal/middleware"
	"gofiber-skeleton/internal/repository"
	"gofiber-skeleton/internal/usecase"
	"gofiber-skeleton/pkg/jwtutil"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database connection
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Initialize JWT utility (Ed25519)
	jwt := jwtutil.NewJWT(cfg.JWTPrivateKey, cfg.JWTPublicKey)

	// Initialize repositories and usecases
	bookingRepo := repository.NewBookingRepository(db)
	bookingUC := usecase.NewBookingUsecase(bookingRepo)

	// Create Fiber app
	app := fiber.New()

	// Apply JWT middleware
	app.Use(middleware.JWTMiddleware(jwt))

	// Register routes
	handler.RegisterRoutes(app, bookingUC)

	// Start server in a goroutine
	go func() {
		if err := app.Listen(":" + cfg.Port); err != nil {
			log.Printf("server stopped: %v", err)
		}
	}()
	log.Printf("server is running on port %s", cfg.Port)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("failed to shut down server gracefully: %v", err)
	}

	log.Println("server exited")
}
