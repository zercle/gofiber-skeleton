package app

import (
	"database/sql"

	validator "github.com/go-playground/validator/v10"
	"github.com/samber/do/v2"
	"github.com/zercle/gofiber-skeleton/internal/ordermodule"
	"github.com/zercle/gofiber-skeleton/internal/productmodule"
	"github.com/zercle/gofiber-skeleton/internal/usermodule"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/config"
	orderhandler "github.com/zercle/gofiber-skeleton/internal/ordermodule/handler"
	orderrepository "github.com/zercle/gofiber-skeleton/internal/ordermodule/repository"
	orderusecase "github.com/zercle/gofiber-skeleton/internal/ordermodule/usecase"
	producthandler "github.com/zercle/gofiber-skeleton/internal/productmodule/handler"
	productrepository "github.com/zercle/gofiber-skeleton/internal/productmodule/repository"
	productusecase "github.com/zercle/gofiber-skeleton/internal/productmodule/usecase"
	userhandler "github.com/zercle/gofiber-skeleton/internal/usermodule/handler"
	userrepository "github.com/zercle/gofiber-skeleton/internal/usermodule/repository"
	userusecase "github.com/zercle/gofiber-skeleton/internal/usermodule/usecase"

	"golang.org/x/crypto/bcrypt"
)

func NewInjector(db *sql.DB, cfg *config.Config) *do.RootScope {
	injector := do.New()

	// Provide the database connection
	do.ProvideValue(injector, db)

	// Provide config
	do.ProvideValue(injector, cfg)

	// Provide repositories
	do.Provide(injector, func(i do.Injector) (usermodule.UserRepository, error) {
		db := do.MustInvoke[*sql.DB](i)
		return userrepository.NewUserRepository(db), nil
	})
	do.Provide(injector, func(i do.Injector) (ordermodule.OrderRepository, error) {
		db := do.MustInvoke[*sql.DB](i)
		return orderrepository.NewOrderRepository(db), nil
	})
	do.Provide(injector, func(i do.Injector) (productmodule.ProductRepository, error) {
		db := do.MustInvoke[*sql.DB](i)
		return productrepository.NewProductRepository(db), nil
	})

	// Provide use cases
	do.Provide(injector, func(i do.Injector) (usermodule.UserUseCase, error) {
		userRepo := do.MustInvoke[usermodule.UserRepository](i)
		cfg := do.MustInvoke[*config.Config](i)
		return userusecase.NewUserUseCase(userRepo, cfg.JWT.Secret, bcrypt.DefaultCost), nil
	})
	do.Provide(injector, func(i do.Injector) (ordermodule.OrderUseCase, error) {
		orderRepo := do.MustInvoke[ordermodule.OrderRepository](i)
		productRepo := do.MustInvoke[productmodule.ProductRepository](i)
		return orderusecase.NewOrderUseCase(orderRepo, productRepo), nil
	})
	do.Provide(injector, func(i do.Injector) (productmodule.ProductUseCase, error) {
		productRepo := do.MustInvoke[productmodule.ProductRepository](i)
		return productusecase.NewProductUseCase(productRepo), nil
	})

	// Provide handlers
	do.Provide(injector, func(i do.Injector) (*userhandler.UserHandler, error) {
		userUseCase := do.MustInvoke[usermodule.UserUseCase](i)
		validator := do.MustInvoke[*validator.Validate](i)
		return userhandler.NewUserHandler(userUseCase, validator), nil
	})
	do.Provide(injector, func(i do.Injector) (*orderhandler.OrderHandler, error) {
		orderUseCase := do.MustInvoke[ordermodule.OrderUseCase](i)
		validator := do.MustInvoke[*validator.Validate](i)
		return orderhandler.NewOrderHandler(orderUseCase, validator), nil
	})
	do.Provide(injector, func(i do.Injector) (*producthandler.ProductHandler, error) {
		productUseCase := do.MustInvoke[productmodule.ProductUseCase](i)
		validator := do.MustInvoke[*validator.Validate](i)
		return producthandler.NewProductHandler(productUseCase, validator), nil
	})

	// Provide validator
	do.Provide(injector, func(i do.Injector) (*validator.Validate, error) {
		return validator.New(), nil
	})

	return injector
}
