package providers

import (
	"go.uber.org/fx"

	authHandlers "github.com/zercle/gofiber-skeleton/pkg/domains/auth/api/handlers"
	authRoutes "github.com/zercle/gofiber-skeleton/pkg/domains/auth/api/routes"
	authRepo "github.com/zercle/gofiber-skeleton/pkg/domains/auth/repositories"
	authUseCase "github.com/zercle/gofiber-skeleton/pkg/domains/auth/biz/usecases"

	postHandlers "github.com/zercle/gofiber-skeleton/pkg/domains/posts/api/handlers"
	postRoutes "github.com/zercle/gofiber-skeleton/pkg/domains/posts/api/routes"
	postRepo "github.com/zercle/gofiber-skeleton/pkg/domains/posts/repositories"
	postUseCase "github.com/zercle/gofiber-skeleton/pkg/domains/posts/biz/usecases"

	"github.com/zercle/gofiber-skeleton/internal/infrastructure/config"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/database"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/middleware"
	infraAuthRepo "github.com/zercle/gofiber-skeleton/pkg/domains/auth/store/repositories"
	infraPostRepo "github.com/zercle/gofiber-skeleton/pkg/domains/posts/store/repositories"
)

// InfrastructureModule provides common infrastructure components.
var InfrastructureModule = fx.Options(
	fx.Provide(config.NewConfig),
	fx.Provide(database.NewDB),
	fx.Provide(middleware.NewLogger),
	fx.Provide(middleware.NewRecover),
	fx.Provide(middleware.NewCORS),
	fx.Provide(middleware.NewAuth),
	fx.Invoke(func(db *database.Database, dbCleanup func()) {
		// This fx.Invoke is primarily to ensure the cleanup function is part of the fx graph.
		// The actual defer call will be handled by fx's lifecycle management for providers.
	}),
)

// AuthModule provides all components for the Auth domain.
var AuthModule = fx.Options(
	fx.Provide(fx.Annotate(infraAuthRepo.NewUserRepository, fx.As(new(authRepo.UserRepository)))),
	fx.Provide(authUseCase.NewAuthUseCase),
	fx.Provide(authHandlers.NewAuthHandler),
	fx.Provide(authRoutes.NewAuthRoutes),
)

// PostModule provides all components for the Post domain.
var PostModule = fx.Options(
	fx.Provide(fx.Annotate(infraPostRepo.NewPostRepository, fx.As(new(postRepo.PostRepository)))),
	fx.Provide(postUseCase.NewPostUseCase),
	fx.Provide(postHandlers.NewPostHandler),
	fx.Provide(postRoutes.NewPostRoutes),
)