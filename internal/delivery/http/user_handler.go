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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "data": err.Error()})
	}

	user, err := h.userUseCase.Register(c.Context(), req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": user})
}

// Login handles user login.
func (h *UserHandler) Login(c *fiber.Ctx) error {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "data": err.Error()})
	}

	token, err := h.userUseCase.Login(c.Context(), req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "data": "Invalid credentials"})
	}

	return c.JSON(fiber.Map{"status": "success", "data": fiber.Map{"token": token}})
}
