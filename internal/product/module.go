package product

import (
	"gofiber-skeleton/internal/infra/config"
	"gofiber-skeleton/internal/product/delivery"
	"gofiber-skeleton/internal/product/repository"
	"gofiber-skeleton/internal/product/usecase"

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
	productRepo := repository.NewProductRepository(m.DB)
	productUsecase := usecase.NewProductUsecase(productRepo)
	productHandler := delivery.NewProductDelivery(productUsecase)

	productGroup := app.Group("/products")
	productGroup.Post("/", productHandler.CreateProduct)
	productGroup.Get("/:id", productHandler.GetProductByID)
	productGroup.Put("/:id", productHandler.UpdateProduct)
	productGroup.Delete("/:id", productHandler.DeleteProduct)
}