//go:generate mockgen -source=user_usecase.go -destination=../mocks/mock_user_usecase.go -package=mocks

package usecase

import (
	"context"
	"gofiber-skeleton/internal/user/domain"
	"gofiber-skeleton/internal/infra/types"
)

type UserUsecase interface {
	GetUser(ctx context.Context, id types.UUIDv7) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) error
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, id types.UUIDv7) error
}