package order

import (
	"gofiber-skeleton/internal/infra/auth"
	"gofiber-skeleton/internal/order/delivery"
	"gofiber-skeleton/internal/order/repository"
	"gofiber-skeleton/internal/order/usecase"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

// RegisterModule initializes and registers the order module for both REST and gRPC.
func RegisterModule(app *fiber.App, grpcServer *grpc.Server, db *gorm.DB, jwtService auth.JWTService) {
	orderRepo := repository.NewOrderRepository(db)
	orderUc := usecase.NewOrderUsecase(orderRepo)

	// Register REST handlers
	delivery.NewOrderHandler(app, orderUc, jwtService)
	// Register gRPC server
	delivery.NewGrpcOrderServer(grpcServer, orderUc)
}
