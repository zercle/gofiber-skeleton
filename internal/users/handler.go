package users

import (
	"github.com/gofiber/fiber/v2"
	helpers "github.com/zercle/gofiber-helpers"
	"github.com/zercle/gofiber-skelton/pkg/domain"
)

type UserHandler struct {
	UserUsecase domain.UserUsecase
}

// NewPostHandler will initialize the post resource endpoint
func InitUserHandler(router fiber.Router, userUsecase domain.UserUsecase) {
	handler := &UserHandler{
		UserUsecase: userUsecase,
	}

	router.Get("/:id", handler.GetUser())
	router.Post("/", handler.CreateUser())
	router.Patch("/:id", handler.EditUser())
	router.Delete("/:id", handler.DeleteUser())
}

func (h *UserHandler) GetUser() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		responseForm := helpers.ResponseForm{}

		userId := c.Params("id")

		user, err := h.UserUsecase.GetUser(userId)

		responseForm.Result = map[string]interface{}{
			"user": user,
		}

		if err == nil {
			responseForm.Success = true
		}
		return c.JSON(responseForm)
	}
}

func (h *UserHandler) CreateUser() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		responseForm := helpers.ResponseForm{}

		if err == nil {
			responseForm.Success = true
		}
		return c.JSON(responseForm)
	}
}

func (h *UserHandler) EditUser() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		responseForm := helpers.ResponseForm{}

		if err == nil {
			responseForm.Success = true
		}
		return c.JSON(responseForm)
	}
}

func (h *UserHandler) DeleteUser() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		responseForm := helpers.ResponseForm{}

		if err == nil {
			responseForm.Success = true
		}
		return c.JSON(responseForm)
	}
}
