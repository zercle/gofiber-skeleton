package delivery

import (
	"context"
	"gofiber-skeleton/api/user"
	
	"gofiber-skeleton/internal/user/usecase"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(app *fiber.App, userUsecase usecase.UserUsecase) {
	h := &UserHandler{
		userUsecase: userUsecase,
	}

	// REST Endpoints
	app.Get("/users/:id", h.GetUserByID)
	app.Post("/users", h.CreateUser)
	app.Put("/users/:id", h.UpdateUser)
	app.Delete("/users/:id", h.DeleteUser)
}

// REST Handlers
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	// TODO: Implement logic to get user by ID
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "GetUserByID"})
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	// TODO: Implement logic to create user
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "CreateUser"})
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	// TODO: Implement logic to update user
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "UpdateUser"})
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	// TODO: Implement logic to delete user
	return c.SendStatus(fiber.StatusNoContent)
}

// gRPC Server
type GrpcUserServer struct {
	user.UnimplementedUserServiceServer
	userUsecase usecase.UserUsecase
}

func NewGrpcUserServer(grpcServer *grpc.Server, userUsecase usecase.UserUsecase) {
	s := &GrpcUserServer{
		userUsecase: userUsecase,
	}
	user.RegisterUserServiceServer(grpcServer, s)
}

// gRPC Handlers
func (s *GrpcUserServer) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.User, error) {
	// TODO: Implement logic to get user
	return &user.User{Id: req.Id, Username: "test", Email: "test@example.com"}, nil
}

func (s *GrpcUserServer) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.User, error) {
	// TODO: Implement logic to create user
	return &user.User{Username: req.Username, Email: req.Email}, nil
}

func (s *GrpcUserServer) UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (*user.User, error) {
	// TODO: Implement logic to update user
	return &user.User{Id: req.Id, Username: req.Username, Email: req.Email}, nil
}

func (s *GrpcUserServer) DeleteUser(ctx context.Context, req *user.DeleteUserRequest) (*user.DeleteUserRequest, error) {
	// TODO: Implement logic to delete user
	return &user.DeleteUserRequest{Id: req.Id}, nil
}
