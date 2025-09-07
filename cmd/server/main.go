package main

import (
	"go.uber.org/fx"

	"github.com/zercle/gofiber-skeleton/internal/infrastructure/config"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/database"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/middleware"
	"github.com/zercle/gofiber-skeleton/internal/shared/container"
)

// @title GoFiber Skeleton API
// @version 1.0
// @description A production-ready Go Fiber backend template with clean architecture
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:3000
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	fx.New(
		// Configuration
		fx.Provide(config.NewConfig),

		// Database
		fx.Provide(database.NewDB),
		fx.Invoke(func(db *database.Database, dbCleanup func()) {
			// This is a placeholder for `dbCleanup()` to be called on program shutdown.
			// The actual defer call will be handled by fx's lifecycle management for providers.
			// This fx.Invoke is primarily to ensure the cleanup function is part of the fx graph.
		}),

		// HTTP Server
		fx.Provide(container.NewFiberApp),

		// Middlewares
		fx.Provide(middleware.NewLogger),
		fx.Provide(middleware.NewRecover),
		fx.Provide(middleware.NewCORS),
		fx.Provide(middleware.NewAuth),

		// Start server
		fx.Invoke(container.RegisterRoutes),
		fx.Invoke(container.StartServer),
	).Run()
}
