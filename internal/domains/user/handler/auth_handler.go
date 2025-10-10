package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/zercle/gofiber-skeleton/internal/domains/user/entity"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/usecase"
	"github.com/zercle/gofiber-skeleton/internal/response"
	"github.com/zercle/gofiber-skeleton/internal/validator"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	authUsecase usecase.AuthUsecase
	validator   *validator.Validator
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
		validator:   validator.New(),
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user account
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body entity.RegisterRequest true "Registration request"
// @Success 201 {object} response.JSendResponse{data=entity.UserResponse}
// @Failure 400 {object} response.JSendResponse
// @Failure 409 {object} response.JSendResponse
// @Failure 500 {object} response.JSendResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req entity.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", nil)
	}

	// Validate request
	if errs := h.validator.Validate(&req); len(errs) > 0 {
		return response.UnprocessableEntity(c, "Validation failed", errs)
	}

	// Register user
	user, err := h.authUsecase.Register(c.Context(), &req)
	if err != nil {
		if contains(err.Error(), "already exists") {
			return response.Conflict(c, err.Error())
		}
		return response.InternalServerError(c, "Failed to register user")
	}

	return response.Created(c, entity.UserResponse{
		User: user.ToPublic(),
	})
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body entity.LoginRequest true "Login request"
// @Success 200 {object} response.JSendResponse{data=entity.LoginResponse}
// @Failure 400 {object} response.JSendResponse
// @Failure 401 {object} response.JSendResponse
// @Failure 500 {object} response.JSendResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req entity.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", nil)
	}

	// Validate request
	if errs := h.validator.Validate(&req); len(errs) > 0 {
		return response.UnprocessableEntity(c, "Validation failed", errs)
	}

	// Login user
	loginResp, err := h.authUsecase.Login(c.Context(), &req)
	if err != nil {
		if contains(err.Error(), "invalid credentials") {
			return response.Unauthorized(c, "Invalid email or password")
		}
		if contains(err.Error(), "deactivated") {
			return response.Forbidden(c, "Account is deactivated")
		}
		return response.InternalServerError(c, "Failed to login")
	}

	return response.Success(c, loginResp)
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get authenticated user's profile
// @Tags User
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.JSendResponse{data=entity.UserResponse}
// @Failure 401 {object} response.JSendResponse
// @Failure 404 {object} response.JSendResponse
// @Failure 500 {object} response.JSendResponse
// @Router /users/me [get]
func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "Unauthorized")
	}

	user, err := h.authUsecase.GetUserByID(c.Context(), userID)
	if err != nil {
		return response.NotFound(c, "User not found")
	}

	return response.Success(c, entity.UserResponse{
		User: user.ToPublic(),
	})
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update authenticated user's profile information
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body entity.UpdateProfileRequest true "Update profile request"
// @Success 200 {object} response.JSendResponse{data=entity.UserResponse}
// @Failure 400 {object} response.JSendResponse
// @Failure 401 {object} response.JSendResponse
// @Failure 422 {object} response.JSendResponse
// @Failure 500 {object} response.JSendResponse
// @Router /users/me [put]
func (h *AuthHandler) UpdateProfile(c *fiber.Ctx) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "Unauthorized")
	}

	var req entity.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", nil)
	}

	// Validate request
	if errs := h.validator.Validate(&req); len(errs) > 0 {
		return response.UnprocessableEntity(c, "Validation failed", errs)
	}

	user, err := h.authUsecase.UpdateProfile(c.Context(), userID, &req)
	if err != nil {
		return response.InternalServerError(c, "Failed to update profile")
	}

	return response.Success(c, entity.UserResponse{
		User: user.ToPublic(),
	})
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change authenticated user's password
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body entity.ChangePasswordRequest true "Change password request"
// @Success 200 {object} response.JSendResponse
// @Failure 400 {object} response.JSendResponse
// @Failure 401 {object} response.JSendResponse
// @Failure 422 {object} response.JSendResponse
// @Failure 500 {object} response.JSendResponse
// @Router /users/me/password [put]
func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "Unauthorized")
	}

	var req entity.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", nil)
	}

	// Validate request
	if errs := h.validator.Validate(&req); len(errs) > 0 {
		return response.UnprocessableEntity(c, "Validation failed", errs)
	}

	if err := h.authUsecase.ChangePassword(c.Context(), userID, &req); err != nil {
		if contains(err.Error(), "invalid old password") {
			return response.BadRequest(c, "Invalid old password", nil)
		}
		return response.InternalServerError(c, "Failed to change password")
	}

	return response.Success(c, fiber.Map{
		"message": "Password changed successfully",
	})
}

// getUserIDFromContext retrieves user ID from Fiber context
func getUserIDFromContext(c *fiber.Ctx) (uuid.UUID, error) {
	userIDStr := c.Locals("userID")
	if userIDStr == nil {
		return uuid.Nil, fiber.NewError(fiber.StatusUnauthorized, "User ID not found in context")
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		return uuid.Nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid user ID format")
	}

	return userID, nil
}

// contains checks if a string contains a substring
func contains(str, substr string) bool {
	return len(str) >= len(substr) && (str == substr || len(substr) == 0 ||
		(len(str) > 0 && len(substr) > 0 && stringContains(str, substr)))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
