package delivery

import (
	"context"
	"gofiber-skeleton/api/user"
	"gofiber-skeleton/pkg/jsend"
	"net/http"

	"gofiber-skeleton/internal/infra/auth"
	"gofiber-skeleton/internal/infra/middleware"
	"gofiber-skeleton/internal/user/usecase"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(app *fiber.App, uc usecase.UserUsecase, jwtService auth.JWTService) {
	handler := &UserHandler{
		userUsecase: uc,
	}

	api := app.Group("/api/v1")

	// Public routes
	api.Post("/register", handler.Register)
	api.Post("/login", handler.Login)

	// Protected routes
	protected := api.Group("/users", middleware.Protected(jwtService))
	protected.Get("/me", handler.GetProfile)
}

// Login godoc
// ... (swagger docs)
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return jsend.Fail(c, map[string]string{"error": "Invalid request body"})
	}

	token, err := h.userUsecase.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return jsend.Error(c, "Invalid credentials", http.StatusUnauthorized, http.StatusUnauthorized)
	}

	return jsend.Success(c, fiber.Map{"token": token})
}

// GetProfile godoc
// ... (swagger docs)
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	// The userID is added to c.Locals by the Protected middleware
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return jsend.Error(c, "Internal server error", http.StatusInternalServerError, http.StatusInternalServerError)
	}

	// user, err := h.userUsecase.GetUserByID(c.Context(), userID)
	// if err != nil {
	// 	return jsend.Error(c, "User not found", http.StatusNotFound, http.StatusNotFound)
	// }

	// For demonstration:
	user := fiber.Map{"id": userID, "email": "user@example.com"}

	return jsend.Success(c, user)
}

// ... other handlers like Register, also using jsend
func (h *UserHandler) Register(c *fiber.Ctx) error {
	// ... parsing and validation logic
	// err := h.userUsecase.CreateUser(...)
	// if err != nil {
	// 	return jsend.Error(c, "Failed to create user", http.StatusBadRequest, http.StatusBadRequest)
	// }
	return jsend.Success(c, fiber.Map{"message": "User created successfully"})
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

func (s *GrpcUserServer) DeleteUser(ctx context.Context, req *user.DeleteUserRequest) (*emptypb.Empty, error) {
	// TODO: Implement logic to delete user
	return &emptypb.Empty{}, nil
}
