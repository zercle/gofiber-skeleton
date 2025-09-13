package main

import (
	"go.uber.org/fx"

	server "github.com/zercle/gofiber-skeleton/internal/app/http"
	"github.com/zercle/gofiber-skeleton/internal/app/providers"
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
		providers.Module,

		// HTTP Server
		fx.Provide(server.NewFiberApp),

		// Start server
		fx.Invoke(server.RegisterRoutes),
		fx.Invoke(server.StartServer),
	).Run()
}
