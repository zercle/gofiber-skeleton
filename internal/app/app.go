package app

import (
	"context"
	"fmt"
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

// Application holds the application's dependencies and state.
type Application struct {
	Config    *config.Config
	DB        *gorm.DB
	JWT       *jwtutil.JWT
	FiberApp  *fiber.App
	BookingUC usecase.BookingUsecase
}

// NewApplication creates and initializes a new Application.
func NewApplication() (*Application, error) {
	cfg := config.LoadConfig()

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	jwtInstance := jwtutil.NewJWT(cfg.JWTPrivateKey, cfg.JWTPublicKey)

	bookingRepo := repository.NewBookingRepository(db)
	bookingUC := usecase.NewBookingUsecase(bookingRepo)

	fiberApp := fiber.New()

	return &Application{
		Config:    cfg,
		DB:        db,
		JWT:       jwtInstance,
		FiberApp:  fiberApp,
		BookingUC: bookingUC,
	}, nil
}

// Start registers routes, middleware, and starts the HTTP server.
func (a *Application) Start() {
	a.FiberApp.Use(middleware.JWTMiddleware(a.JWT))
	handler.RegisterRoutes(a.FiberApp, a.BookingUC)

	go func() {
		if err := a.FiberApp.Listen(":" + a.Config.Port); err != nil {
			// Log non-nil errors, especially if server stops unexpectedly.
			// http.ErrServerClosed is a normal error on graceful shutdown.
			if err.Error() != "http: Server closed" {
				log.Printf("server error: %v", err)
			}
		}
	}()
	log.Printf("server is running on port %s", a.Config.Port)
}

// WaitForShutdown listens for OS signals and handles graceful shutdown.
func (a *Application) WaitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := a.FiberApp.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("failed to shut down server gracefully: %v", err)
	}

	log.Println("server exited")
}