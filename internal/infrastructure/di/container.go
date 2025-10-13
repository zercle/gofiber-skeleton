package di

import (
	"context"
	"log"

	"github.com/samber/do"
	"github.com/zercle/gofiber-skeleton/internal/domains/post/handler"
	postRepo "github.com/zercle/gofiber-skeleton/internal/domains/post/repository"
	postUsecase "github.com/zercle/gofiber-skeleton/internal/domains/post/usecase"
	userHandler "github.com/zercle/gofiber-skeleton/internal/domains/user/handler"
	userRepo "github.com/zercle/gofiber-skeleton/internal/domains/user/repository"
	userUsecase "github.com/zercle/gofiber-skeleton/internal/domains/user/usecase"
	userRouter "github.com/zercle/gofiber-skeleton/internal/domains/user/router"
	postRouter "github.com/zercle/gofiber-skeleton/internal/domains/post/router"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/config"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/database"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/database/sqlc"
	"github.com/zercle/gofiber-skeleton/internal/shared/middleware"
	sharedRouter "github.com/zercle/gofiber-skeleton/internal/shared/router"
)

func NewServiceContainer(ctx context.Context) *do.Injector {
	injector := do.New()

	// Inject config
	do.Provide(injector, func(i *do.Injector) (*config.Config, error) {
		cfg, err := config.Load()
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}
		return cfg, nil
	})

	// Inject database
	do.Provide(injector, func(i *do.Injector) (*database.Database, error) {
		cfg := do.MustInvoke[*config.Config](i)
		db, err := database.NewDatabase(
			ctx,
			cfg.GetDSN(),
			cfg.Database.MaxOpen,
			cfg.Database.MaxIdle,
		)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		return db, nil
	})

	// Inject SQLC Querier
	do.Provide(injector, func(i *do.Injector) (sqlc.Querier, error) {
		db := do.MustInvoke[*database.Database](i)
		return sqlc.New(db.GetPool()), nil
	})

	// Inject User Repository
	do.Provide(injector, func(i *do.Injector) (userRepo.UserRepository, error) {
		queries := do.MustInvoke[sqlc.Querier](i)
		return userRepo.NewUserRepository(queries), nil
	})

	// Inject User Usecase
	do.Provide(injector, func(i *do.Injector) (userUsecase.UserUsecase, error) {
		userRepo := do.MustInvoke[userRepo.UserRepository](i)
		cfg := do.MustInvoke[*config.Config](i)
		return userUsecase.NewUserUsecase(userRepo, cfg.JWT.Secret), nil
	})

	// Inject User Handler
	do.Provide(injector, func(i *do.Injector) (*userHandler.UserHandler, error) {
		userUsecase := do.MustInvoke[userUsecase.UserUsecase](i)
		return userHandler.NewUserHandler(userUsecase), nil
	})

	// Inject Post Repository
	do.Provide(injector, func(i *do.Injector) (postRepo.PostRepository, error) {
		queries := do.MustInvoke[sqlc.Querier](i)
		return postRepo.NewPostRepository(queries), nil
	})

	// Inject Post Usecase
	do.Provide(injector, func(i *do.Injector) (postUsecase.PostUsecase, error) {
		postRepo := do.MustInvoke[postRepo.PostRepository](i)
		return postUsecase.NewPostUsecase(postRepo), nil
	})

	// Inject Post Handler
	do.Provide(injector, func(i *do.Injector) (*handler.PostHandler, error) {
		postUsecase := do.MustInvoke[postUsecase.PostUsecase](i)
		return handler.NewPostHandler(postUsecase), nil
	})

	// Inject Auth Middleware
	do.Provide(injector, func(i *do.Injector) (*middleware.AuthMiddleware, error) {
		cfg := do.MustInvoke[*config.Config](i)
		return middleware.NewAuthMiddleware(cfg.JWT.Secret), nil
	})

	// Inject User Router
	do.Provide(injector, userRouter.NewUserRouter)

	// Inject Post Router
	do.Provide(injector, postRouter.NewPostRouter)

	// Inject shared routers
	do.Provide(injector, func(i *do.Injector) ([]sharedRouter.Router, error) {
		userRouter := do.MustInvoke[*userRouter.UserRouter](i)
		postRouter := do.MustInvoke[*postRouter.PostRouter](i)
		return []sharedRouter.Router{
			userRouter,
			postRouter,
		}, nil
	})

	return injector
}