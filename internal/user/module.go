package user

import (
	"gofiber-skeleton/internal/infra/auth"
	"gofiber-skeleton/internal/user/delivery"
	"gofiber-skeleton/internal/user/repository"
	"gofiber-skeleton/internal/user/usecase"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

// RegisterModule initializes and registers the user module for both REST and gRPC.
func RegisterModule(app *fiber.App, grpcServer *grpc.Server, db *gorm.DB, jwtService auth.JWTService) {
	userRepo := repository.NewUserRepository(db)
	userUc := usecase.NewUserUsecase(userRepo, jwtService)

	// Register REST handlers
	delivery.NewUserHandler(app, userUc, jwtService)
	// Register gRPC server
	delivery.NewGrpcUserServer(grpcServer, userUc)
}
