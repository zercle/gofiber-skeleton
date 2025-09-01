package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/zercle/gofiber-skeleton/internal/auth"
	"github.com/zercle/gofiber-skeleton/internal/core/domain"
)

type AuthHandler struct {
	authUsecase auth.AuthUsecase
	validate    *validator.Validate
}

func NewAuthHandler(authUsecase auth.AuthUsecase, validate *validator.Validate) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
		validate:    validate,
	}
}

func (h *AuthHandler) Login(ctx *fiber.Ctx) error {
	var req auth.LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return domain.NewError(domain.ErrCodeBadRequest, "invalid request body")
	}

	if err := h.validate.Struct(&req); err != nil {
		return domain.NewError(domain.ErrCodeBadRequest, "invalid request body")
	}

	res, err := h.authUsecase.Login(ctx, req.Email, req.Password)
	if err != nil {
		return err
	}

	return ctx.JSON(res)
}
