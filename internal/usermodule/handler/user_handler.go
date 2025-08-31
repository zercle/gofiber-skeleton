package userhandler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/zercle/gofiber-skeleton/internal/usermodule"
	"github.com/zercle/gofiber-skeleton/pkg/jsend"
)

type UserHandler struct {
	userUseCase usermodule.UserUseCase
	validator   *validator.Validate
}

// NewUserHandler creates a new user handler
func NewUserHandler(userUseCase usermodule.UserUseCase, validator *validator.Validate) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
		validator:   validator,
	}
}

// RegisterRequest represents the request body for user registration
type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role"` // Role can be empty, not required
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Register handles user registration
func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return jsend.Fail(c, jsend.Empty, "Invalid request body")
	}

	if err := h.validator.Struct(&req); err != nil {
		return jsend.Fail(c, jsend.Empty, err.Error())
	}

	// Create user
	user, err := h.userUseCase.Register(req.Username, req.Password, req.Role)
	if err != nil {
		if err.Error() == "username already exists" {
			return jsend.Fail(c, jsend.Empty, err.Error())
		}
		return jsend.Error(c, err.Error(), 0, http.StatusInternalServerError)
	}

	return jsend.SuccessWithStatus(c, fiber.Map{
		"message": "User registered successfully",
		"user":    user,
	}, http.StatusCreated)
}

// Login handles user login
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return jsend.Fail(c, jsend.Empty, "Invalid request body")
	}

		if err := h.validator.Struct(&req); err != nil {
			return jsend.Fail(c, jsend.Empty, err.Error())
		}
	// Login user
	token, user, err := h.userUseCase.Login(req.Username, req.Password)
	if err != nil {
		if err.Error() == "invalid credentials" {
			return jsend.Error(c, err.Error(), 0, http.StatusUnauthorized)
		}
		return jsend.Error(c, err.Error(), 0, http.StatusInternalServerError)
	}

	return jsend.Success(c, fiber.Map{
		"message": "Login successful",
		"token":   token,
		"user":    user,
	})
}

// GetByID handles getting a user by ID
func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return jsend.Fail(c, jsend.Empty, "User ID is required")
	}

	user, err := h.userUseCase.GetByID(id)
	if err != nil {
		return jsend.Fail(c, jsend.Empty, "User not found")
	}

	return jsend.Success(c, fiber.Map{
		"user": user,
	})
}

// UpdateRole handles updating a user's role
func (h *UserHandler) UpdateRole(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return jsend.Fail(c, jsend.Empty, "User ID is required")
	}

	var req struct {
		Role string `json:"role" validate:"required,oneof=admin customer"`
	}
	if err := c.BodyParser(&req); err != nil {
		return jsend.Fail(c, jsend.Empty, "Invalid request body")
	}

	if err := h.validator.Struct(&req); err != nil {
		return jsend.Fail(c, jsend.Empty, err.Error())
	}

	if err := h.userUseCase.UpdateRole(id, req.Role); err != nil {
		return jsend.Error(c, err.Error(), 0, http.StatusInternalServerError)
	}

	return jsend.Success(c, fiber.Map{
		"message": "User role updated successfully",
	})
}
