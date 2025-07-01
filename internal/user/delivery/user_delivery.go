package delivery

import (
	"context"
	"gofiber-skeleton/api/user"

	"gofiber-skeleton/internal/user/usecase"
	"gofiber-skeleton/pkg/jsend"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserDelivery struct {
	userUsecase usecase.UserUsecase
}

func NewUserDelivery(userUsecase usecase.UserUsecase) *UserDelivery {
	return &UserDelivery{
		userUsecase: userUsecase,
	}
}

// REST Handlers
// GetUserByID godoc
// @Summary Get user by ID
// @Description Get user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} jsend.JSendResponse{data=domain.User} "Success"
// @Failure 400 {object} jsend.JSendResponse{data=string} "Bad Request"
// @Failure 404 {object} jsend.JSendResponse{data=string} "Not Found"
// @Failure 500 {object} jsend.JSendResponse{data=string} "Internal Server Error"
// @Router /users/{id} [get]
func (d *UserDelivery) GetUserByID(c *fiber.Ctx) error {
	// TODO: Implement logic to get user by ID
	return jsend.Success(c, fiber.Map{"message": "GetUserByID"})
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body domain.User true "User object"
// @Success 201 {object} jsend.JSendResponse{data=domain.User} "Created"
// @Failure 400 {object} jsend.JSendResponse{data=string} "Bad Request"
// @Failure 500 {object} jsend.JSendResponse{data=string} "Internal Server Error"
// @Router /users [post]
func (d *UserDelivery) CreateUser(c *fiber.Ctx) error {
	// TODO: Implement logic to create user
	return jsend.Success(c, fiber.Map{"message": "CreateUser"})
}

// UpdateUser godoc
// @Summary Update an existing user
// @Description Update an existing user
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body domain.User true "User object"
// @Success 200 {object} jsend.JSendResponse{data=domain.User} "Success"
// @Failure 400 {object} jsend.JSendResponse{data=string} "Bad Request"
// @Failure 404 {object} jsend.JSendResponse{data=string} "Not Found"
// @Failure 500 {object} jsend.JSendResponse{data=string} "Internal Server Error"
// @Router /users/{id} [put]
func (d *UserDelivery) UpdateUser(c *fiber.Ctx) error {
	// TODO: Implement logic to update user
	return jsend.Success(c, fiber.Map{"message": "UpdateUser"})
}

// DeleteUser godoc
// @Summary Delete a user by ID
// @Description Delete a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} jsend.JSendResponse{data=string} "Success"
// @Failure 400 {object} jsend.JSendResponse{data=string} "Bad Request"
// @Failure 404 {object} jsend.JSendResponse{data=string} "Not Found"
// @Failure 500 {object} jsend.JSendResponse{data=string} "Internal Server Error"
// @Router /users/{id} [delete]
func (d *UserDelivery) DeleteUser(c *fiber.Ctx) error {
	// TODO: Implement logic to delete user
	return jsend.Success(c, fiber.Map{"message": "DeleteUser"})
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

func (s *GrpcUserServer) DeleteUser(ctx context.Context, req *user.DeleteUserRequest) (*emptypb.Empty, error) {
	// TODO: Implement logic to delete user
	return &emptypb.Empty{}, nil
}