package delivery

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/samber/do/v2"

	"github.com/zercle/gofiber-skeleton/internal/domains/user/entity"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/usecase"
	"github.com/zercle/gofiber-skeleton/internal/middleware"
	"github.com/zercle/gofiber-skeleton/pkg/response"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	userUsecase usecase.UserUsecase
}

// NewUserHandler creates a new user handler
func NewUserHandler(injector do.Injector) (*UserHandler, error) {
	userUsecase := do.MustInvoke[usecase.UserUsecase](injector)

	return &UserHandler{
		userUsecase: userUsecase,
	}, nil
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user with email, password and full name
// @Tags users
// @Accept json
// @Produce json
// @Param request body entity.CreateUserRequest true "Registration request"
// @Success 201 {object} response.Response{data=entity.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/auth/register [post]
func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req entity.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
	}

	user, err := h.userUsecase.Register(c.Context(), &req)
	if err != nil {
		if err == entity.ErrEmailExists {
			return response.Error(c, fiber.StatusConflict, "EMAIL_EXISTS", err.Error())
		}
		if domainErr, ok := err.(*entity.DomainError); ok {
			return response.Error(c, fiber.StatusBadRequest, domainErr.Code, domainErr.Message)
		}
		return response.Error(c, fiber.StatusInternalServerError, "REGISTRATION_FAILED", "Failed to register user")
	}

	return response.Success(c, fiber.StatusCreated, "User registered successfully", user)
}

// Login handles user authentication
// @Summary Login user
// @Description Authenticate user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body entity.LoginRequest true "Login request"
// @Success 200 {object} response.Response{data=entity.LoginResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/auth/login [post]
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req entity.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
	}

	result, err := h.userUsecase.Login(c.Context(), &req)
	if err != nil {
		if err == entity.ErrInvalidPassword || err == entity.ErrUserNotFound || err == entity.ErrUserInactive {
			return response.Error(c, fiber.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password")
		}
		if domainErr, ok := err.(*entity.DomainError); ok {
			return response.Error(c, fiber.StatusBadRequest, domainErr.Code, domainErr.Message)
		}
		return response.Error(c, fiber.StatusInternalServerError, "LOGIN_FAILED", "Failed to authenticate user")
	}

	return response.Success(c, fiber.StatusOK, "Login successful", result)
}

// GetProfile handles getting user profile
// @Summary Get user profile
// @Description Get authenticated user's profile
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=entity.UserResponse}
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/users/profile [get]
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userIDStr := middleware.GetUserID(c)
	if userIDStr == "" {
		return response.Error(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "INVALID_USER_ID", "Invalid user ID")
	}

	user, err := h.userUsecase.GetProfile(c.Context(), userID)
	if err != nil {
		if err == entity.ErrUserNotFound {
			return response.Error(c, fiber.StatusNotFound, "USER_NOT_FOUND", "User not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, "PROFILE_FAILED", "Failed to get user profile")
	}

	return response.Success(c, fiber.StatusOK, "Profile retrieved successfully", user)
}

// UpdateProfile handles updating user profile
// @Summary Update user profile
// @Description Update authenticated user's profile
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body entity.UpdateUserRequest true "Update profile request"
// @Success 200 {object} response.Response{data=entity.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 409 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/users/profile [put]
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userIDStr := middleware.GetUserID(c)
	if userIDStr == "" {
		return response.Error(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "INVALID_USER_ID", "Invalid user ID")
	}

	var req entity.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
	}

	user, err := h.userUsecase.UpdateProfile(c.Context(), userID, &req)
	if err != nil {
		if err == entity.ErrUserNotFound {
			return response.Error(c, fiber.StatusNotFound, "USER_NOT_FOUND", "User not found")
		}
		if err == entity.ErrEmailExists {
			return response.Error(c, fiber.StatusConflict, "EMAIL_EXISTS", "Email already exists")
		}
		if domainErr, ok := err.(*entity.DomainError); ok {
			return response.Error(c, fiber.StatusBadRequest, domainErr.Code, domainErr.Message)
		}
		return response.Error(c, fiber.StatusInternalServerError, "UPDATE_FAILED", "Failed to update profile")
	}

	return response.Success(c, fiber.StatusOK, "Profile updated successfully", user)
}

// ChangePassword handles changing user password
// @Summary Change user password
// @Description Change authenticated user's password
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body entity.ChangePasswordRequest true "Change password request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/users/change-password [post]
func (h *UserHandler) ChangePassword(c *fiber.Ctx) error {
	userIDStr := middleware.GetUserID(c)
	if userIDStr == "" {
		return response.Error(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "INVALID_USER_ID", "Invalid user ID")
	}

	var req entity.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
	}

	err = h.userUsecase.ChangePassword(c.Context(), userID, &req)
	if err != nil {
		if err == entity.ErrUserNotFound {
			return response.Error(c, fiber.StatusNotFound, "USER_NOT_FOUND", "User not found")
		}
		if err == entity.ErrInvalidPassword {
			return response.Error(c, fiber.StatusBadRequest, "INVALID_CURRENT_PASSWORD", "Current password is incorrect")
		}
		if domainErr, ok := err.(*entity.DomainError); ok {
			return response.Error(c, fiber.StatusBadRequest, domainErr.Code, domainErr.Message)
		}
		return response.Error(c, fiber.StatusInternalServerError, "PASSWORD_CHANGE_FAILED", "Failed to change password")
	}

	return response.Success(c, fiber.StatusOK, "Password changed successfully", nil)
}

// ListUsers handles listing users (admin only)
// @Summary List users
// @Description List all users with pagination (admin only)
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.Response{data=[]entity.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/users [get]
func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	// Parse pagination parameters
	page := 1
	limit := 10

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	offset := (page - 1) * limit

	users, err := h.userUsecase.ListUsers(c.Context(), limit, offset)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "LIST_USERS_FAILED", "Failed to list users")
	}

	return response.Success(c, fiber.StatusOK, "Users retrieved successfully", users)
}

// DeactivateUser handles deactivating a user (admin only)
// @Summary Deactivate user
// @Description Deactivate a user account (admin only)
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) DeactivateUser(c *fiber.Ctx) error {
	userIDStr := c.Params("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "INVALID_USER_ID", "Invalid user ID")
	}

	err = h.userUsecase.DeactivateUser(c.Context(), userID)
	if err != nil {
		if err == entity.ErrUserNotFound {
			return response.Error(c, fiber.StatusNotFound, "USER_NOT_FOUND", "User not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, "DEACTIVATE_FAILED", "Failed to deactivate user")
	}

	return response.Success(c, fiber.StatusOK, "User deactivated successfully", nil)
}