package http

import (
	"errors"
	"gofiber-skeleton/internal/user/usecase"
	"gofiber-skeleton/pkg/constant"

	"github.com/gofiber/fiber/v2"
)

// UserHandler handles HTTP requests related to User operations.
type UserHandler struct {
	userUseCase usecase.UserUseCase
}

// NewHTTPUserHandler creates a new instance of UserHandler with the provided User use case.
func NewHTTPUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase: userUseCase}
}

// RegisterRequest represents the expected JSON payload for user registration.
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role" validate:"required,oneof=admin customer"`
}

// Register handles user registration.
// @Summary Register a new user
// @Description Registers a user with a username, password, and role.
// @Tags Users
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "User registration payload"
// @Success 201 {object} map[string]interface{} "User registered successfully message"
// @Failure 400 {object} map[string]string "Invalid input or username already exists"
// @Failure 500 {object} map[string]string "Failed to register user"
// @Router /api/users/register [post]
func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	if req.Username == "" || req.Password == "" || req.Role == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Username, password, and role cannot be empty",
		})
	}

	// Set default role to customer if not specified
	if req.Role == "" {
		req.Role = "customer"
	}

	_, err := h.userUseCase.Register(c.Context(), req.Username, req.Password, req.Role)
	if err != nil {
		if errors.Is(err, constant.ErrUsernameAlreadyExists) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Username already exists",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to register user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "User registered successfully",
	})
}

// LoginRequest represents the expected JSON payload for user login.
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Login handles user login and returns an authentication token.
// @Summary User login
// @Description Authenticates a user and returns a JWT token upon success.
// @Tags Users
// @Accept json
// @Produce json
// @Param request body LoginRequest true "User login payload"
// @Success 200 {object} map[string]interface{} "Authentication token"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Router /api/users/login [post]
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	if req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Username and password cannot be empty",
		})
	}

	token, err := h.userUseCase.Login(c.Context(), req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid credentials",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   fiber.Map{"token": token},
	})
}
