package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/zercle/gofiber-skeleton/internal/domains/user/entity"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/usecase"
	"github.com/zercle/gofiber-skeleton/internal/shared/middleware"
	"github.com/zercle/gofiber-skeleton/internal/shared/response"
	"github.com/zercle/gofiber-skeleton/internal/shared/validator"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param request body entity.CreateUserRequest true "User registration data"
// @Success 201 {object} response.Response{data=entity.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/auth/register [post]
func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req entity.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	if err := validator.Validate(&req); err != nil {
		return response.ValidationError(c, err.Error())
	}

	user, err := h.userUsecase.Register(c.Context(), &req)
	if err != nil {
		switch err {
		case usecase.ErrEmailAlreadyUsed:
			return response.Conflict(c, "Email already used", err.Error())
		default:
			return response.InternalServerError(c, "Failed to register user", err.Error())
		}
	}

	return response.Created(c, "User registered successfully", user.ToResponse())
}

// Login godoc
// @Summary User login
// @Description Authenticate a user and return a JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param request body entity.LoginRequest true "User login data"
// @Success 200 {object} response.Response{data=entity.LoginResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/auth/login [post]
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req entity.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	if err := validator.Validate(&req); err != nil {
		return response.ValidationError(c, err.Error())
	}

	loginResp, err := h.userUsecase.Login(c.Context(), &req)
	if err != nil {
		switch err {
		case usecase.ErrUserNotFound, usecase.ErrInvalidPassword:
			return response.Unauthorized(c, "Invalid credentials", err.Error())
		default:
			return response.InternalServerError(c, "Failed to login", err.Error())
		}
	}

	return response.OK(c, "Login successful", &entity.LoginResponse{
		User:  loginResp.User,
		Token: loginResp.Token,
	})
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get the current user's profile
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=entity.UserResponse}
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/users/profile [get]
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "Invalid user context", err.Error())
	}

	user, err := h.userUsecase.GetByID(c.Context(), userID)
	if err != nil {
		switch err {
		case usecase.ErrUserNotFound:
			return response.NotFound(c, "User not found", err.Error())
		default:
			return response.InternalServerError(c, "Failed to get user", err.Error())
		}
	}

	return response.OK(c, "User profile retrieved successfully", user.ToResponse())
}

// GetByID godoc
// @Summary Get user by ID
// @Description Get a user by their ID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Security BearerAuth
// @Success 200 {object} response.Response{data=entity.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", err.Error())
	}

	user, err := h.userUsecase.GetByID(c.Context(), userID)
	if err != nil {
		switch err {
		case usecase.ErrUserNotFound:
			return response.NotFound(c, "User not found", err.Error())
		default:
			return response.InternalServerError(c, "Failed to get user", err.Error())
		}
	}

	return response.OK(c, "User retrieved successfully", user.ToResponse())
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update the current user's profile
// @Tags users
// @Accept json
// @Produce json
// @Param request body entity.UpdateUserRequest true "User update data"
// @Security BearerAuth
// @Success 200 {object} response.Response{data=entity.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 409 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/users/profile [put]
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "Invalid user context", err.Error())
	}

	var req entity.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	if err := validator.Validate(&req); err != nil {
		return response.ValidationError(c, err.Error())
	}

	user, err := h.userUsecase.Update(c.Context(), userID, &req)
	if err != nil {
		switch err {
		case usecase.ErrUserNotFound:
			return response.NotFound(c, "User not found", err.Error())
		case usecase.ErrEmailAlreadyUsed:
			return response.Conflict(c, "Email already used", err.Error())
		default:
			return response.InternalServerError(c, "Failed to update user", err.Error())
		}
	}

	return response.OK(c, "User updated successfully", user.ToResponse())
}

// DeleteAccount godoc
// @Summary Delete user account
// @Description Delete the current user's account
// @Tags users
// @Security BearerAuth
// @Success 204 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/users/profile [delete]
func (h *UserHandler) DeleteAccount(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return response.Unauthorized(c, "Invalid user context", err.Error())
	}

	if err := h.userUsecase.Delete(c.Context(), userID); err != nil {
		switch err {
		case usecase.ErrUserNotFound:
			return response.NotFound(c, "User not found", err.Error())
		default:
			return response.InternalServerError(c, "Failed to delete user", err.Error())
		}
	}

	return response.NoContent(c, "User deleted successfully")
}

// List godoc
// @Summary List users
// @Description Get a paginated list of users
// @Tags users
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]entity.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/users [get]
func (h *UserHandler) List(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	users, err := h.userUsecase.List(c.Context(), limit, offset)
	if err != nil {
		return response.InternalServerError(c, "Failed to list users", err.Error())
	}

	userResponses := make([]*entity.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = user.ToResponse()
	}

	return response.OK(c, "Users retrieved successfully", userResponses)
}