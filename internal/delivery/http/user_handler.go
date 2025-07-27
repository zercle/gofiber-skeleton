package http

import (
	"gofiber-skeleton/internal/usecases"

	"github.com/gofiber/fiber/v2"
)

// NewUserHandler creates a new UserHandler.
func NewUserHandler(userUseCase usecases.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase: userUseCase}
}

// UserHandler handles HTTP requests for users.
type UserHandler struct {
	userUseCase usecases.UserUseCase
}

// Register handles user registration.
func (h *UserHandler) Register(c *fiber.Ctx) error {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
	}

	if req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Username and password cannot be empty"})
	}

	_, err := h.userUseCase.Register(c.Context(), req.Username, req.Password)
	if err != nil {
		if err.Error() == "username already exists" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Username already exists"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to register user"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "message": "User registered successfully"})
}

// Login handles user login.
func (h *UserHandler) Login(c *fiber.Ctx) error {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
	}

	if req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Username and password cannot be empty"})
	}

	token, err := h.userUseCase.Login(c.Context(), req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid credentials"})
	}

	return c.JSON(fiber.Map{"status": "success", "data": fiber.Map{"token": token}})
}
