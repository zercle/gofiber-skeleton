package app

import (
	"database/sql"

	validator "github.com/go-playground/validator/v10"
	do_v2 "github.com/samber/do/v2"
	"github.com/zercle/gofiber-skeleton/internal/domain"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/config"
	orderhandler "github.com/zercle/gofiber-skeleton/internal/order/handler"
	orderrepository "github.com/zercle/gofiber-skeleton/internal/order/repository"
	orderusecase "github.com/zercle/gofiber-skeleton/internal/order/usecase"
	producthandler "github.com/zercle/gofiber-skeleton/internal/product/handler"
	productrepository "github.com/zercle/gofiber-skeleton/internal/product/repository"
	productusecase "github.com/zercle/gofiber-skeleton/internal/product/usecase"
	userhandler "github.com/zercle/gofiber-skeleton/internal/user/handler"
	userrepository "github.com/zercle/gofiber-skeleton/internal/user/repository"
	userusecase "github.com/zercle/gofiber-skeleton/internal/user/usecase"

	"golang.org/x/crypto/bcrypt"
)

func NewInjector(db *sql.DB, cfg *config.Config) *do_v2.RootScope {
	injector := do_v2.New()

	// Provide the database connection
	do_v2.ProvideValue(injector, db)

	// Provide config
	do_v2.ProvideValue(injector, cfg)

	// Provide repositories
	do_v2.Provide(injector, func(i do_v2.Injector) (domain.UserRepository, error) {
		db := do_v2.MustInvoke[*sql.DB](i)
		return userrepository.NewUserRepository(db), nil
	})
	do_v2.Provide(injector, func(i do_v2.Injector) (domain.OrderRepository, error) {
		db := do_v2.MustInvoke[*sql.DB](i)
		return orderrepository.NewOrderRepository(db), nil
	})
	do_v2.Provide(injector, func(i do_v2.Injector) (domain.ProductRepository, error) {
		db := do_v2.MustInvoke[*sql.DB](i)
		return productrepository.NewProductRepository(db), nil
	})

	// Provide use cases
	do_v2.Provide(injector, func(i do_v2.Injector) (domain.UserUseCase, error) {
		userRepo := do_v2.MustInvoke[domain.UserRepository](i)
		cfg := do_v2.MustInvoke[*config.Config](i)
		return userusecase.NewUserUseCase(userRepo, cfg.JWT.Secret, bcrypt.DefaultCost), nil
	})
	do_v2.Provide(injector, func(i do_v2.Injector) (domain.OrderUseCase, error) {
		orderRepo := do_v2.MustInvoke[domain.OrderRepository](i)
		productRepo := do_v2.MustInvoke[domain.ProductRepository](i)
		return orderusecase.NewOrderUseCase(orderRepo, productRepo), nil
	})
	do_v2.Provide(injector, func(i do_v2.Injector) (domain.ProductUseCase, error) {
		productRepo := do_v2.MustInvoke[domain.ProductRepository](i)
		return productusecase.NewProductUseCase(productRepo), nil
	})

	// Provide handlers
	do_v2.Provide(injector, func(i do_v2.Injector) (*userhandler.UserHandler, error) {
		userUseCase := do_v2.MustInvoke[domain.UserUseCase](i)
		return userhandler.NewUserHandler(userUseCase), nil
	})
	do_v2.Provide(injector, func(i do_v2.Injector) (*orderhandler.OrderHandler, error) {
		orderUseCase := do_v2.MustInvoke[domain.OrderUseCase](i)
		return orderhandler.NewOrderHandler(orderUseCase), nil
	})
	do_v2.Provide(injector, func(i do_v2.Injector) (*producthandler.ProductHandler, error) {
		productUseCase := do_v2.MustInvoke[domain.ProductUseCase](i)
		validator := do_v2.MustInvoke[*validator.Validate](i)
		return producthandler.NewProductHandler(productUseCase, validator), nil
	})

	// Provide validator
	do_v2.Provide(injector, func(i do_v2.Injector) (*validator.Validate, error) {
		return validator.New(), nil
	})

	return injector
}
