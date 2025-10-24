package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zercle/template-go-fiber/internal/domains"
	"github.com/zercle/template-go-fiber/internal/errors"
	"github.com/zercle/template-go-fiber/pkg/response"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	usecase domains.UserUsecase
}

// NewUserHandler creates a new user handler
func NewUserHandler(usecase domains.UserUsecase) *UserHandler {
	return &UserHandler{
		usecase: usecase,
	}
}

// CreateUserRequest is the request body for user creation
type CreateUserRequest struct {
	Email     string  `json:"email" binding:"required,email" example:"user@example.com"`
	Password  string  `json:"password" binding:"required,min=8" example:"password123"`
	FirstName *string `json:"first_name,omitempty" example:"John"`
	LastName  *string `json:"last_name,omitempty" example:"Doe"`
}

// UpdateUserRequest is the request body for user update
type UpdateUserRequest struct {
	Email     *string `json:"email,omitempty" binding:"omitempty,email"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	IsActive  *bool   `json:"is_active,omitempty"`
}

// UserResponse is the response body for user data
type UserResponse struct {
	ID        string  `json:"id"`
	Email     string  `json:"email"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	IsActive  bool    `json:"is_active"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

// RegisterUser registers a new user
// @Summary Register a new user
// @Description Create a new user account with email and password
// @Tags Users
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "User registration request"
// @Success 201 {object} response.Response[UserResponse]
// @Failure 400 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Router /users/register [post]
func (h *UserHandler) RegisterUser(c *fiber.Ctx) error {
	var req CreateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return response.SendError(c, errors.NewValidationError("invalid request body", err))
	}

	input := &domains.RegisterUserInput{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	user, err := h.usecase.RegisterUser(c.Context(), input)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			return response.SendError(c, apiErr)
		}
		return response.SendUnknownError(c, err)
	}

	return response.SendCreated(c, userToResponse(user))
}

// GetUser retrieves a user by ID
// @Summary Get user by ID
// @Description Retrieve a user's information by their ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} response.Response[UserResponse]
// @Failure 404 {object} response.ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := h.usecase.GetUserByID(c.Context(), id)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			return response.SendError(c, apiErr)
		}
		return response.SendUnknownError(c, err)
	}

	return response.SendOK(c, userToResponse(user))
}

// GetUserByEmail retrieves a user by email
// @Summary Get user by email
// @Description Retrieve a user's information by their email
// @Tags Users
// @Accept json
// @Produce json
// @Param email query string true "User email"
// @Success 200 {object} response.Response[UserResponse]
// @Failure 404 {object} response.ErrorResponse
// @Router /users/email [get]
func (h *UserHandler) GetUserByEmail(c *fiber.Ctx) error {
	email := c.Query("email")

	if email == "" {
		return response.SendError(c, errors.NewValidationError("email query parameter is required", nil))
	}

	user, err := h.usecase.GetUserByEmail(c.Context(), email)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			return response.SendError(c, apiErr)
		}
		return response.SendUnknownError(c, err)
	}

	return response.SendOK(c, userToResponse(user))
}

// UpdateUser updates a user
// @Summary Update user information
// @Description Update an existing user's information
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body UpdateUserRequest true "User update request"
// @Success 200 {object} response.Response[UserResponse]
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Security Bearer
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")

	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return response.SendError(c, errors.NewValidationError("invalid request body", err))
	}

	input := &domains.UpdateUserInput{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  req.IsActive,
	}

	user, err := h.usecase.UpdateUser(c.Context(), id, input)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			return response.SendError(c, apiErr)
		}
		return response.SendUnknownError(c, err)
	}

	return response.SendOK(c, userToResponse(user))
}

// DeleteUser deletes a user
// @Summary Delete a user
// @Description Delete an existing user account
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 204
// @Failure 404 {object} response.ErrorResponse
// @Security Bearer
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	err := h.usecase.DeleteUser(c.Context(), id)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			return response.SendError(c, apiErr)
		}
		return response.SendUnknownError(c, err)
	}

	return response.SendNoContent(c)
}

// ListUsers lists all users with pagination
// @Summary List users
// @Description Retrieve a paginated list of users
// @Tags Users
// @Accept json
// @Produce json
// @Param limit query int false "Number of users to return (default 10, max 100)" default(10)
// @Param offset query int false "Number of users to skip" default(0)
// @Success 200 {object} response.Response[[]UserResponse]
// @Router /users [get]
func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)

	users, err := h.usecase.ListUsers(c.Context(), int32(limit), int32(offset))
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			return response.SendError(c, apiErr)
		}
		return response.SendUnknownError(c, err)
	}

	responses := make([]UserResponse, len(users))
	for i, user := range users {
		resp := userToResponse(user)
		responses[i] = resp
	}

	return response.SendOK(c, responses)
}

// Helper function to convert domain user to response
func userToResponse(user *domains.User) UserResponse {
	createdAt := ""
	if !user.CreatedAt.IsZero() {
		createdAt = user.CreatedAt.Format("2006-01-02T15:04:05Z")
	}

	updatedAt := ""
	if !user.UpdatedAt.IsZero() {
		updatedAt = user.UpdatedAt.Format("2006-01-02T15:04:05Z")
	}

	return UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsActive:  user.IsActive,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
