package container

import (
	"github.com/samber/do/v2"
	"github.com/zercle/gofiber-skeleton/internal/config"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/database"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/sqlc"
	"github.com/zercle/gofiber-skeleton/internal/user"
	"github.com/zercle/gofiber-skeleton/internal/user/handler"
	"github.com/zercle/gofiber-skeleton/internal/user/repository"
	"github.com/zercle/gofiber-skeleton/internal/user/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/zercle/gofiber-skeleton/internal/auth"
	authHandler "github.com/zercle/gofiber-skeleton/internal/auth/handler"
	authUsecase "github.com/zercle/gofiber-skeleton/internal/auth/usecase"
)

// SetupContainer initializes and configures the dependency injection container.
// It registers all the application's dependencies, such as configuration,
// database connections, repositories, use cases, and handlers.
// This function is called at the start of the application to build the
// dependency graph that will be used throughout the application's lifecycle.
func SetupContainer() (*do.RootScope, error) {
	container := do.New()

	// Register configuration
	do.ProvideValue(container, config.Load())

	// Register database
	do.Provide(container, func(i do.Injector) (*database.Database, error) {
		cfg := do.MustInvoke[*config.Config](i)
		return database.NewPostgresConnection(cfg)
	})

	// Register SQLC queries
	do.Provide(container, func(i do.Injector) (*sqlc.Queries, error) {
		db := do.MustInvoke[*database.Database](i)
		return sqlc.New(db.DB), nil
	})

	// Register User domain
	do.Provide(container, func(i do.Injector) (user.UserRepository, error) {
		queries := do.MustInvoke[*sqlc.Queries](i)
		return repository.NewUserRepository(queries), nil
	})

	do.Provide(container, func(i do.Injector) (user.UserUsecase, error) {
		userRepo := do.MustInvoke[user.UserRepository](i)
		return usecase.NewUserUsecase(userRepo), nil
	})

	do.Provide(container, func(i do.Injector) (*handler.UserHandler, error) {
		userUsecase := do.MustInvoke[user.UserUsecase](i)
		return handler.NewUserHandler(userUsecase), nil
	})

	// Register Auth domain
	do.Provide(container, func(i do.Injector) (auth.AuthUsecase, error) {
		userUsecase := do.MustInvoke[user.UserUsecase](i)
		cfg := do.MustInvoke[*config.Config](i)
		return authUsecase.NewAuthUsecase(userUsecase, cfg), nil
	})

	do.Provide(container, func(i do.Injector) (*authHandler.AuthHandler, error) {
		authUsecase := do.MustInvoke[auth.AuthUsecase](i)
		validate := do.MustInvoke[*validator.Validate](i)
		return authHandler.NewAuthHandler(authUsecase, validate), nil
	})

	// Register validator
	do.Provide(container, func(i do.Injector) (*validator.Validate, error) {
		return validator.New(), nil
	})

	return container, nil
}