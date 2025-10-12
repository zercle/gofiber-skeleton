package container

import (
	"github.com/samber/do/v2"
	"github.com/zercle/gofiber-skeleton/internal/config"
	"github.com/zercle/gofiber-skeleton/internal/domains/post/delivery"
	postrepository "github.com/zercle/gofiber-skeleton/internal/domains/post/repository"
	"github.com/zercle/gofiber-skeleton/internal/domains/post/usecase"
	userdelivery "github.com/zercle/gofiber-skeleton/internal/domains/user/delivery"
	userrepository "github.com/zercle/gofiber-skeleton/internal/domains/user/repository"
	userusecase "github.com/zercle/gofiber-skeleton/internal/domains/user/usecase"
	"github.com/zercle/gofiber-skeleton/pkg/database"
)

func NewContainer(cfg *config.Config) do.Injector {
	container := do.New()

	// Configuration
	do.ProvideValue(container, cfg)

	// Database
	do.Provide(container, database.NewDatabase)

	// User Domain
	do.Provide(container, userrepository.NewUserRepository)
	do.Provide(container, userusecase.NewUserUsecase)
	do.Provide(container, userdelivery.NewUserHandler)

	// Post Domain
	do.Provide(container, postrepository.NewPostRepository)
	do.Provide(container, usecase.NewPostUsecase)
	do.Provide(container, delivery.NewPostHandler)

	return container
}
