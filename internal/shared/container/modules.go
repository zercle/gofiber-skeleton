package container

import (
	"go.uber.org/fx"

	"github.com/zercle/gofiber-skeleton/internal/infrastructure/config"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/database"

	authRepo "github.com/zercle/gofiber-skeleton/internal/domains/auth/repositories"
	authUseCase "github.com/zercle/gofiber-skeleton/internal/domains/auth/usecases"
	postRepo "github.com/zercle/gofiber-skeleton/internal/domains/posts/repositories"
	postUseCase "github.com/zercle/gofiber-skeleton/internal/domains/posts/usecases"

	infraAuthRepo "github.com/zercle/gofiber-skeleton/internal/infrastructure/repositories/auth"
	infraPostRepo "github.com/zercle/gofiber-skeleton/internal/infrastructure/repositories/posts"
)

// InfrastructureModule provides common infrastructure components.
var InfrastructureModule = fx.Options(
	fx.Provide(config.NewConfig),
	fx.Provide(database.NewDB),
	fx.Invoke(func(db *database.Database, dbCleanup func()) {
		// This fx.Invoke is primarily to ensure the cleanup function is part of the fx graph.
		// The actual defer call will be handled by fx's lifecycle management for providers.
	}),
)

// AuthModule provides all components for the Auth domain.
var AuthModule = fx.Options(
	fx.Provide(fx.Annotate(infraAuthRepo.NewUserRepository, fx.As(new(authRepo.UserRepository)))),
	fx.Provide(authUseCase.NewAuthUseCase),
	// fx.Provide(authHandler.NewAuthHandler), // Handlers will be added later
)

// PostModule provides all components for the Post domain.
var PostModule = fx.Options(
	fx.Provide(fx.Annotate(infraPostRepo.NewPostRepository, fx.As(new(postRepo.PostRepository)))),
	fx.Provide(postUseCase.NewPostUseCase),
	// fx.Provide(postHandler.NewPostHandler), // Handlers will be added later
)
