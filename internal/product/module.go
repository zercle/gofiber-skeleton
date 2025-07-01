package product

import (
	"gofiber-skeleton/internal/infra/auth"
	"gofiber-skeleton/internal/product/delivery"
	"gofiber-skeleton/internal/product/repository"
	"gofiber-skeleton/internal/product/usecase"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

// RegisterModule initializes and registers the product module for both REST and gRPC.
func RegisterModule(app *fiber.App, grpcServer *grpc.Server, db *gorm.DB, jwtService auth.JWTService) {
	productRepo := repository.NewProductRepository(db)
	productUc := usecase.NewProductUsecase(productRepo)

	// Register REST handlers
	delivery.NewProductHandler(app, productUc, jwtService)
	// Register gRPC server
	delivery.NewGrpcProductServer(grpcServer, productUc)
}
