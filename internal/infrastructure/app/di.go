package app

import (
	"database/sql"

	validator "github.com/go-playground/validator/v10"
	"github.com/samber/do/v2"
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

	demohandler "github.com/zercle/gofiber-skeleton/internal/demo/handler"
	demorepository "github.com/zercle/gofiber-skeleton/internal/demo/repository"
	demousecase "github.com/zercle/gofiber-skeleton/internal/demo/usecase"
	"golang.org/x/crypto/bcrypt"
)

func NewInjector(db *sql.DB, cfg *config.Config) *do.RootScope {
	injector := do.New()

	// Provide the database connection
	do.ProvideValue(injector, db)

	// Provide config
	do.ProvideValue(injector, cfg)

	// Provide repositories
	do.Provide(injector, func(i do.Injector) (domain.UserRepository, error) {
		db := do.MustInvoke[*sql.DB](i)
		return userrepository.NewUserRepository(db), nil
	})
	do.Provide(injector, func(i do.Injector) (domain.OrderRepository, error) {
		db := do.MustInvoke[*sql.DB](i)
		return orderrepository.NewOrderRepository(db), nil
	})
	do.Provide(injector, func(i do.Injector) (domain.ProductRepository, error) {
		db := do.MustInvoke[*sql.DB](i)
		return productrepository.NewProductRepository(db), nil
	})
	do.Provide(injector, func(i do.Injector) (demorepository.DemoRepository, error) {
		db := do.MustInvoke[*sql.DB](i)
		return demorepository.NewDemoRepository(db), nil
	})

	// Provide use cases
	do.Provide(injector, func(i do.Injector) (domain.UserUseCase, error) {
		userRepo := do.MustInvoke[domain.UserRepository](i)
		cfg := do.MustInvoke[*config.Config](i)
		return userusecase.NewUserUseCase(userRepo, cfg.JWT.Secret, bcrypt.DefaultCost), nil
	})
	do.Provide(injector, func(i do.Injector) (domain.OrderUseCase, error) {
		orderRepo := do.MustInvoke[domain.OrderRepository](i)
		productRepo := do.MustInvoke[domain.ProductRepository](i)
		return orderusecase.NewOrderUseCase(orderRepo, productRepo), nil
	})
	do.Provide(injector, func(i do.Injector) (domain.ProductUseCase, error) {
		productRepo := do.MustInvoke[domain.ProductRepository](i)
		return productusecase.NewProductUseCase(productRepo), nil
	})
	do.Provide(injector, func(i do.Injector) (demousecase.DemoUseCase, error) {
		demoRepo := do.MustInvoke[demorepository.DemoRepository](i)
		return demousecase.NewDemoUseCase(demoRepo), nil
	})

	// Provide handlers
	do.Provide(injector, func(i do.Injector) (*userhandler.UserHandler, error) {
		userUseCase := do.MustInvoke[domain.UserUseCase](i)
		return userhandler.NewUserHandler(userUseCase), nil
	})
	do.Provide(injector, func(i do.Injector) (*orderhandler.OrderHandler, error) {
		orderUseCase := do.MustInvoke[domain.OrderUseCase](i)
		return orderhandler.NewOrderHandler(orderUseCase), nil
	})
	do.Provide(injector, func(i do.Injector) (*producthandler.ProductHandler, error) {
		productUseCase := do.MustInvoke[domain.ProductUseCase](i)
		validator := do.MustInvoke[*validator.Validate](i)
		return producthandler.NewProductHandler(productUseCase, validator), nil
	})
	do.Provide(injector, func(i do.Injector) (*demohandler.DemoHandler, error) {
		demoUseCase := do.MustInvoke[demousecase.DemoUseCase](i)
		cfg := do.MustInvoke[*config.Config](i)
		return demohandler.NewDemoHandler(cfg, demoUseCase), nil
	})

	// Provide validator
	do.Provide(injector, func(i do.Injector) (*validator.Validate, error) {
		return validator.New(), nil
	})

	return injector
}
