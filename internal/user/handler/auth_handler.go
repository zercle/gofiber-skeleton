package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/zercle/gofiber-skeleton/internal/logger"
	"github.com/zercle/gofiber-skeleton/internal/response"
	"github.com/zercle/gofiber-skeleton/internal/user/usecase"
	"github.com/zercle/gofiber-skeleton/pkg/validator"
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

// Register godoc
// @Summary Register a new user
// @Description Create a new user account with username, email, and password
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body registerRequest true "Registration details"
// @Success 201 {object} response.JSendResponse "User registered successfully"
// @Failure 400 {object} response.JSendResponse "Invalid request"
// @Failure 500 {object} response.JSendResponse "Internal server error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	req := new(registerRequest)
	if err := c.BodyParser(req); err != nil {
		logger.GetLogger().Error().Err(err).Msg("Failed to parse register request")
		return response.Fail(c, http.StatusBadRequest, fiber.Map{"error": "Invalid request body"})
	}

	if err := validator.ValidateRequest(c, req); err != nil {
		return err
	}

	user, err := h.authUsecase.Register(c.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		logger.GetLogger().Error().Err(err).Msg("Failed to register user")
		return response.Error(c, http.StatusInternalServerError, "Failed to register user", 1001)
	}

	return response.Success(c, http.StatusCreated, fiber.Map{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body loginRequest true "Login credentials"
// @Success 200 {object} response.JSendResponse "Login successful with JWT token"
// @Failure 400 {object} response.JSendResponse "Invalid request"
// @Failure 401 {object} response.JSendResponse "Invalid credentials"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	req := new(loginRequest)
	if err := c.BodyParser(req); err != nil {
		logger.GetLogger().Error().Err(err).Msg("Failed to parse login request")
		return response.Fail(c, http.StatusBadRequest, fiber.Map{"error": "Invalid request body"})
	}

	if err := validator.ValidateRequest(c, req); err != nil {
		return err
	}

	token, user, err := h.authUsecase.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		logger.GetLogger().Error().Err(err).Msg("Failed to login user")
		return response.Fail(c, http.StatusUnauthorized, fiber.Map{"error": "Invalid credentials"})
	}

	return response.Success(c, http.StatusOK, fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
	})
}

func RegisterAuthRoutes(router fiber.Router, authUsecase usecase.AuthUsecase) {
	handler := NewAuthHandler(authUsecase)
	router.Post("/register", handler.Register)
	router.Post("/login", handler.Login)
}
