//go:generate mockgen -source=user_usecase.go -destination=../mocks/mock_user_usecase.go -package=mocks

package usecase

import (
	"context"
	"gofiber-skeleton/internal/user/domain"
)

type UserUsecase interface {
	GetUser(ctx context.Context, id uint) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) error
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, id uint) error
}