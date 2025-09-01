package handler

import (
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/validator"
	"github.com/zercle/gofiber-skeleton/internal/pkg/response"
	"github.com/zercle/gofiber-skeleton/internal/user"
)

var registerPayloadPool = sync.Pool{
	New: func() any {
		return new(user.RegisterPayload)
	},
}

// UserHandler handles user-related HTTP requests.
type UserHandler struct {
	userUsecase user.UserUsecase
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(userUsecase user.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

// Register handles user registration.
// @Summary Register a new user
// @Description Create a new user account with email and password
// @Tags users
// @Accept json
// @Produce json
// @Param user body user.RegisterPayload true "User registration data"
// @Success 201 {object} response.JSendResponse{data=object{user=user.User}} "User created successfully"
// @Failure 400 {object} response.JSendResponse{data=object{body_parser=string}} "Invalid request body"
// @Failure 409 {object} response.JSendResponse{message=string} "User with this email already exists"
// @Failure 422 {object} response.JSendResponse{data=object{email=string,password=string}} "Validation errors"
// @Failure 500 {object} response.JSendResponse{message=string} "Internal server error"
// @Router /api/v1/users/register [post]
func (h *UserHandler) Register(c *fiber.Ctx) error {
	payload := registerPayloadPool.Get().(*user.RegisterPayload)
	defer registerPayloadPool.Put(payload)
	// It's good practice to reset the payload before reusing it.
	payload.Email = ""
	payload.Password = ""

	if err := c.BodyParser(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Fail(map[string]any{
			"body_parser": "invalid request body: " + err.Error(),
		}))
	}

	if validationErrors := validator.ValidateStruct(payload); len(validationErrors) > 0 {
		return validator.HandleValidationErrors(c, validationErrors)
	}

	createdUser, err := h.userUsecase.Register(*payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(err.Error()))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(map[string]any{
		"user": createdUser,
	}))
}

// GetAll handles fetching all users.
func (h *UserHandler) GetAll(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Not Implemented")
}

// GetByID handles fetching a user by ID.
func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Not Implemented")
}

// Update handles updating a user.
func (h *UserHandler) Update(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Not Implemented")
}

// Delete handles deleting a user.
func (h *UserHandler) Delete(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("Not Implemented")
}
