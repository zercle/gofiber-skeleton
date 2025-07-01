package app

import (
	"gofiber-skeleton/internal/infra/config"
	"gofiber-skeleton/internal/order"
	"gofiber-skeleton/internal/product"
	"gofiber-skeleton/internal/user"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// App represents the main application structure.
type App struct {
	FiberApp *fiber.App
	DB       *gorm.DB
	Config   *config.Config
}

// NewApp creates a new instance of the App.
func NewApp(fiberApp *fiber.App, db *gorm.DB, cfg *config.Config) *App {
	return &App{
		FiberApp: fiberApp,
		DB:       db,
		Config:   cfg,
	}
}

// SetupRoutes sets up all the application routes.
func (a *App) SetupRoutes() {
	// Initialize modules
	userModule := user.NewModule(a.DB, a.Config)
	productModule := product.NewModule(a.DB, a.Config)
	orderModule := order.NewModule(a.DB, a.Config)

	// Setup routes for each module
	userModule.SetupRoutes(a.FiberApp)
	productModule.SetupRoutes(a.FiberApp)
	orderModule.SetupRoutes(a.FiberApp)
}