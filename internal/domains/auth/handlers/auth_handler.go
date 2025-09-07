package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/zercle/gofiber-skeleton/internal/domains/auth/models"
	"github.com/zercle/gofiber-skeleton/internal/domains/auth/usecases"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/validation"
	"github.com/zercle/gofiber-skeleton/internal/shared/jsend"
)

type AuthHandler struct {
	authUseCase usecases.AuthUseCase
	validator   *validator.Validate
}

func NewAuthHandler(authUseCase usecases.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
		validator:   validation.NewValidator(),
	}
}

// Login godoc
// @Summary User login
// @Description Authenticate user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login credentials"
// @Success 200 {object} jsend.JSendResponse{data=models.LoginResponse}
// @Failure 400 {object} jsend.JSendResponse
// @Failure 401 {object} jsend.JSendResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return jsend.SendBadRequest(c, map[string]string{
			"message": "Invalid request body",
		})
	}

	if err := h.validator.Struct(&req); err != nil {
		var errors []jsend.ValidationError
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, jsend.ValidationError{
				Field:   err.Field(),
				Message: err.Tag(),
			})
		}
		return jsend.SendValidationError(c, errors)
	}

	response, err := h.authUseCase.Login(c.Context(), &req)
	if err != nil {
		return jsend.SendUnauthorized(c, map[string]string{
			"message": err.Error(),
		})
	}

	return jsend.SendSuccess(c, response)
}

// Register godoc
// @Summary User registration
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Registration data"
// @Success 201 {object} jsend.JSendResponse{data=models.RegisterResponse}
// @Failure 400 {object} jsend.JSendResponse
// @Failure 409 {object} jsend.JSendResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return jsend.SendBadRequest(c, map[string]string{
			"message": "Invalid request body",
		})
	}

	if err := h.validator.Struct(&req); err != nil {
		var errors []jsend.ValidationError
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, jsend.ValidationError{
				Field:   err.Field(),
				Message: err.Tag(),
			})
		}
		return jsend.SendValidationError(c, errors)
	}

	response, err := h.authUseCase.Register(c.Context(), &req)
	if err != nil {
		if err.Error() == "email already registered" {
			return jsend.SendFail(c, fiber.StatusConflict, map[string]string{
				"message": err.Error(),
			})
		}
		return jsend.SendInternalServerError(c, err.Error())
	}

	return jsend.SendSuccess(c.Status(fiber.StatusCreated), response)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Generate new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} jsend.JSendResponse{data=models.RefreshTokenResponse}
// @Failure 400 {object} jsend.JSendResponse
// @Failure 401 {object} jsend.JSendResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req models.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return jsend.SendBadRequest(c, map[string]string{
			"message": "Invalid request body",
		})
	}

	if err := h.validator.Struct(&req); err != nil {
		var errors []jsend.ValidationError
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, jsend.ValidationError{
				Field:   err.Field(),
				Message: err.Tag(),
			})
		}
		return jsend.SendValidationError(c, errors)
	}

	response, err := h.authUseCase.RefreshToken(c.Context(), &req)
	if err != nil {
		return jsend.SendUnauthorized(c, map[string]string{
			"message": err.Error(),
		})
	}

	return jsend.SendSuccess(c, response)
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get authenticated user's profile information
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} jsend.JSendResponse{data=models.UserData}
// @Failure 401 {object} jsend.JSendResponse
// @Failure 404 {object} jsend.JSendResponse
// @Router /auth/profile [get]
func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	userIDStr := c.Locals("userID").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return jsend.SendBadRequest(c, map[string]string{
			"message": "Invalid user ID",
		})
	}

	profile, err := h.authUseCase.GetProfile(c.Context(), userID)
	if err != nil {
		return jsend.SendNotFound(c, map[string]string{
			"message": err.Error(),
		})
	}

	return jsend.SendSuccess(c, profile)
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change the authenticated user's password
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.ChangePasswordRequest true "Password change data"
// @Success 204
// @Failure 400 {object} jsend.JSendResponse
// @Failure 401 {object} jsend.JSendResponse
// @Router /auth/change-password [put]
func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	userIDStr := c.Locals("userID").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return jsend.SendBadRequest(c, map[string]string{
			"message": "Invalid user ID",
		})
	}

	var req models.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return jsend.SendBadRequest(c, map[string]string{
			"message": "Invalid request body",
		})
	}

	if err := h.validator.Struct(&req); err != nil {
		var errors []jsend.ValidationError
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, jsend.ValidationError{
				Field:   err.Field(),
				Message: err.Tag(),
			})
		}
		return jsend.SendValidationError(c, errors)
	}

	err = h.authUseCase.ChangePassword(c.Context(), userID, &req)
	if err != nil {
		if err.Error() == "current password is incorrect" {
			return jsend.SendBadRequest(c, map[string]string{
				"message": err.Error(),
			})
		}
		return jsend.SendInternalServerError(c, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
