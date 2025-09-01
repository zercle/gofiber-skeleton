package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zercle/gofiber-skeleton/internal/core/domain"
)

type AuthUsecase interface {
	Login(ctx *fiber.Ctx, email, password string) (*LoginResponse, *domain.Error)
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type Claims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
}