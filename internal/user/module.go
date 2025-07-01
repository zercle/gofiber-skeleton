package user

import (
	"gofiber-skeleton/internal/infra/config"
	"gofiber-skeleton/internal/user/delivery"
	"gofiber-skeleton/internal/user/repository"
	"gofiber-skeleton/internal/user/usecase"

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
	userRepo := repository.NewUserRepository(m.DB)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := delivery.NewUserDelivery(userUsecase)

	userGroup := app.Group("/users")
	userGroup.Post("/", userHandler.CreateUser)
	userGroup.Get("/:id", userHandler.GetUserByID)
	userGroup.Put("/:id", userHandler.UpdateUser)
	userGroup.Delete("/:id", userHandler.DeleteUser)
}