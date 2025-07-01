package order

import (
	"gofiber-skeleton/internal/infra/config"
	"gofiber-skeleton/internal/order/delivery"
	"gofiber-skeleton/internal/order/repository"
	"gofiber-skeleton/internal/order/usecase"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Module struct {
	DB     *gorm.DB
	Config *config.Config
}

func NewModule(db *gorm.DB, cfg *config.Config) *Module {
	return &Module{
		DB:     db,
		Config: cfg,
	}
}

func (m *Module) SetupRoutes(app *fiber.App) {
	orderRepo := repository.NewOrderRepository(m.DB)
	orderUsecase := usecase.NewOrderUsecase(orderRepo)
	orderHandler := delivery.NewOrderDelivery(orderUsecase)

	orderGroup := app.Group("/orders")
	orderGroup.Post("/", orderHandler.CreateOrder)
	orderGroup.Get("/:id", orderHandler.GetOrderByID)
	orderGroup.Put("/:id", orderHandler.UpdateOrder)
	orderGroup.Delete("/:id", orderHandler.DeleteOrder)
}