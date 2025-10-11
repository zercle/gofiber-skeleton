package delivery

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/entity"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/usecase"
	"github.com/zercle/gofiber-skeleton/internal/middleware"
	"github.com/zercle/gofiber-skeleton/pkg/response"
)

// UserHandler handles HTTP requests for user domain
type UserHandler struct {
	userUsecase usecase.UserUsecase
}

// NewUserHandler creates a new user handler
func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body entity.RegisterRequest true "Registration request"
// @Success 201 {object} response.Response{data=entity.User}
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /auth/register [post]
func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req entity.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	user, err := h.userUsecase.Register(c.Context(), &req)
	if err != nil {
		if errors.Is(err, errors.New("email already registered")) ||
			errors.Is(err, errors.New("username already taken")) {
			return response.Conflict(c, err.Error())
		}
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, user, "User registered successfully")
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return access token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body entity.LoginRequest true "Login request"
// @Success 200 {object} response.Response{data=entity.LoginResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /auth/login [post]
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req entity.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	loginResp, err := h.userUsecase.Login(c.Context(), &req)
	if err != nil {
		return response.Unauthorized(c, err.Error())
	}

	return response.Success(c, loginResp, "Login successful")
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get the current user's profile
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=entity.User}
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /users/profile [get]
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userIDStr := middleware.GetUserID(c)
	if userIDStr == "" {
		return response.Unauthorized(c, "User not authenticated")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID")
	}

	user, err := h.userUsecase.GetProfile(c.Context(), userID)
	if err != nil {
		return response.NotFound(c, "User")
	}

	return response.Success(c, user)
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Get a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} response.Response{data=entity.User}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid user ID")
	}

	user, err := h.userUsecase.GetByID(c.Context(), userID)
	if err != nil {
		return response.NotFound(c, "User")
	}

	return response.Success(c, user)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update the current user's profile
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body entity.UpdateUserRequest true "Update request"
// @Success 200 {object} response.Response{data=entity.User}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 409 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /users/profile [put]
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userIDStr := middleware.GetUserID(c)
	if userIDStr == "" {
		return response.Unauthorized(c, "User not authenticated")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID")
	}

	var req entity.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	user, err := h.userUsecase.Update(c.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, errors.New("email already registered")) ||
			errors.Is(err, errors.New("username already taken")) {
			return response.Conflict(c, err.Error())
		}
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, user, "Profile updated successfully")
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change the current user's password
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body entity.ChangePasswordRequest true "Change password request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /users/password [put]
func (h *UserHandler) ChangePassword(c *fiber.Ctx) error {
	userIDStr := middleware.GetUserID(c)
	if userIDStr == "" {
		return response.Unauthorized(c, "User not authenticated")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID")
	}

	var req entity.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	if err := h.userUsecase.ChangePassword(c.Context(), userID, &req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, nil, "Password changed successfully")
}

// ListUsers godoc
// @Summary List users
// @Description Get a paginated list of users
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(20)
// @Success 200 {object} response.PaginatedResponse{data=[]entity.User}
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /users [get]
func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 20)

	userList, err := h.userUsecase.List(c.Context(), page, perPage)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	totalPages := int(userList.TotalCount) / perPage
	if int(userList.TotalCount)%perPage != 0 {
		totalPages++
	}

	return response.Paginated(c, userList.Users, response.PaginationMeta{
		CurrentPage: page,
		PerPage:     perPage,
		TotalPages:  totalPages,
		TotalCount:  userList.TotalCount,
	})
}

// DeactivateUser godoc
// @Summary Deactivate user
// @Description Deactivate a user account
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /users/{id}/deactivate [post]
func (h *UserHandler) DeactivateUser(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid user ID")
	}

	if err := h.userUsecase.Deactivate(c.Context(), userID); err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, nil, "User deactivated successfully")
}

// ActivateUser godoc
// @Summary Activate user
// @Description Activate a user account
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /users/{id}/activate [post]
func (h *UserHandler) ActivateUser(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid user ID")
	}

	if err := h.userUsecase.Activate(c.Context(), userID); err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, nil, "User activated successfully")
}

// DeleteUser godoc
// @Summary Delete user
// @Description Permanently delete a user account
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid user ID")
	}

	if err := h.userUsecase.Delete(c.Context(), userID); err != nil {
		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, nil, "User deleted successfully")
}
