package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/zercle/gofiber-skeleton/internal/domains/auth/models"
	"github.com/zercle/gofiber-skeleton/internal/domains/auth/usecases"
	"github.com/zercle/gofiber-skeleton/internal/shared/types"
	"github.com/zercle/gofiber-skeleton/pkg/utils"
)

type AuthHandler struct {
	authUsecase usecases.AuthUsecase
}

func NewAuthHandler(authUsecase usecases.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
	}
}

// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Register request"
// @Success 201 {object} models.AuthResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendFail(c, fiber.StatusBadRequest, map[string]interface{}{
			"message": "Invalid request body",
		})
	}

	if validationErrors := utils.ValidateStruct(req); validationErrors != nil {
		return utils.SendValidationError(c, validationErrors)
	}

	result, err := h.authUsecase.Register(c.Context(), req)
	if err != nil {
		switch err {
		case types.ErrEmailExists:
			return utils.SendConflict(c, "Email already exists")
		default:
			return utils.SendInternalError(c, "Failed to register user")
		}
	}

	return utils.SendCreated(c, result)
}

// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login request"
// @Success 200 {object} models.AuthResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendFail(c, fiber.StatusBadRequest, map[string]interface{}{
			"message": "Invalid request body",
		})
	}

	if validationErrors := utils.ValidateStruct(req); validationErrors != nil {
		return utils.SendValidationError(c, validationErrors)
	}

	result, err := h.authUsecase.Login(c.Context(), req)
	if err != nil {
		switch err {
		case types.ErrInvalidCredentials:
			return utils.SendUnauthorized(c, "Invalid credentials")
		case types.ErrUnauthorized:
			return utils.SendUnauthorized(c, "Account is inactive")
		default:
			return utils.SendInternalError(c, "Failed to login")
		}
	}

	return utils.SendSuccess(c, result)
}

// @Summary Get user profile
// @Description Get current user's profile information
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.UserResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/profile [get]
func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	result, err := h.authUsecase.GetProfile(c.Context(), userID)
	if err != nil {
		switch err {
		case types.ErrUserNotFound:
			return utils.SendNotFound(c, "User not found")
		default:
			return utils.SendInternalError(c, "Failed to get profile")
		}
	}

	return utils.SendSuccess(c, result)
}

// @Summary Update user profile
// @Description Update current user's profile information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.UpdateProfileRequest true "Update profile request"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/profile [put]
func (h *AuthHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var req models.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendFail(c, fiber.StatusBadRequest, map[string]interface{}{
			"message": "Invalid request body",
		})
	}

	if validationErrors := utils.ValidateStruct(req); validationErrors != nil {
		return utils.SendValidationError(c, validationErrors)
	}

	result, err := h.authUsecase.UpdateProfile(c.Context(), userID, req)
	if err != nil {
		switch err {
		case types.ErrUserNotFound:
			return utils.SendNotFound(c, "User not found")
		default:
			return utils.SendInternalError(c, "Failed to update profile")
		}
	}

	return utils.SendSuccess(c, result)
}

// @Summary Change password
// @Description Change current user's password
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.ChangePasswordRequest true "Change password request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var req models.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendFail(c, fiber.StatusBadRequest, map[string]interface{}{
			"message": "Invalid request body",
		})
	}

	if validationErrors := utils.ValidateStruct(req); validationErrors != nil {
		return utils.SendValidationError(c, validationErrors)
	}

	err := h.authUsecase.ChangePassword(c.Context(), userID, req)
	if err != nil {
		switch err {
		case types.ErrUserNotFound:
			return utils.SendNotFound(c, "User not found")
		case types.ErrInvalidCredentials:
			return utils.SendUnauthorized(c, "Current password is incorrect")
		default:
			return utils.SendInternalError(c, "Failed to change password")
		}
	}

	return utils.SendSuccess(c, map[string]interface{}{
		"message": "Password changed successfully",
	})
}

// @Summary List users
// @Description Get list of all users (admin only)
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} models.UsersListResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/users [get]
func (h *AuthHandler) ListUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	result, err := h.authUsecase.ListUsers(c.Context(), page, pageSize)
	if err != nil {
		return utils.SendInternalError(c, "Failed to list users")
	}

	return utils.SendSuccess(c, result)
}

// @Summary Get user by ID
// @Description Get user information by ID (admin only)
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} models.UserResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/users/{id} [get]
func (h *AuthHandler) GetUser(c *fiber.Ctx) error {
	userID := c.Params("id")

	result, err := h.authUsecase.GetUser(c.Context(), userID)
	if err != nil {
		switch err {
		case types.ErrUserNotFound:
			return utils.SendNotFound(c, "User not found")
		default:
			return utils.SendInternalError(c, "Failed to get user")
		}
	}

	return utils.SendSuccess(c, result)
}

// @Summary Deactivate user
// @Description Deactivate user account (admin only)
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/users/{id}/deactivate [post]
func (h *AuthHandler) DeactivateUser(c *fiber.Ctx) error {
	userID := c.Params("id")

	err := h.authUsecase.DeactivateUser(c.Context(), userID)
	if err != nil {
		switch err {
		case types.ErrUserNotFound:
			return utils.SendNotFound(c, "User not found")
		default:
			return utils.SendInternalError(c, "Failed to deactivate user")
		}
	}

	return utils.SendSuccess(c, map[string]interface{}{
		"message": "User deactivated successfully",
	})
}

// @Summary Activate user
// @Description Activate user account (admin only)
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/users/{id}/activate [post]
func (h *AuthHandler) ActivateUser(c *fiber.Ctx) error {
	userID := c.Params("id")

	err := h.authUsecase.ActivateUser(c.Context(), userID)
	if err != nil {
		switch err {
		case types.ErrUserNotFound:
			return utils.SendNotFound(c, "User not found")
		default:
			return utils.SendInternalError(c, "Failed to activate user")
		}
	}

	return utils.SendSuccess(c, map[string]interface{}{
		"message": "User activated successfully",
	})
}

// @Summary Delete user
// @Description Delete user account (admin only)
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/users/{id} [delete]
func (h *AuthHandler) DeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")

	err := h.authUsecase.DeleteUser(c.Context(), userID)
	if err != nil {
		switch err {
		case types.ErrUserNotFound:
			return utils.SendNotFound(c, "User not found")
		default:
			return utils.SendInternalError(c, "Failed to delete user")
		}
	}

	return utils.SendSuccess(c, map[string]interface{}{
		"message": "User deleted successfully",
	})
}