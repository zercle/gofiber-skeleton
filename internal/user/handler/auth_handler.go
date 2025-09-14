package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/zercle/gofiber-skeleton/internal/logger"
	"github.com/zercle/gofiber-skeleton/internal/user/usecase"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

type registerRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	req := new(registerRequest)
	if err := c.BodyParser(req); err != nil {
		logger.GetLogger().Error().Err(err).Msg("Failed to parse register request")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}

	// TODO: Add validation using a library like go-playground/validator

	user, err := h.authUsecase.Register(c.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		logger.GetLogger().Error().Err(err).Msg("Failed to register user")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to register user"})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"id":        user.ID,
		"username":  user.Username,
		"email":     user.Email,
		"createdAt": user.CreatedAt,
		"updatedAt": user.UpdatedAt,
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	req := new(loginRequest)
	if err := c.BodyParser(req); err != nil {
		logger.GetLogger().Error().Err(err).Msg("Failed to parse login request")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	}

	// TODO: Add validation using a library like go-playground/validator

	token, user, err := h.authUsecase.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		logger.GetLogger().Error().Err(err).Msg("Failed to login user")
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid credentials"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":        user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"createdAt": user.CreatedAt,
			"updatedAt": user.UpdatedAt,
		},
	})
}

func RegisterAuthRoutes(router fiber.Router, authUsecase usecase.AuthUsecase) {
	handler := NewAuthHandler(authUsecase)
	router.Post("/register", handler.Register)
	router.Post("/login", handler.Login)
}