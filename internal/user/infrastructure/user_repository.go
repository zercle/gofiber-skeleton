//go:generate mockgen -source=user_repository.go -destination=../mocks/mock_user_repository.go -package=mocks
package infrastructure

import (
	"context"
	"gofiber-skeleton/internal/user/domain"
)

type UserRepository interface {
	GetUser(ctx context.Context, id uint) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) error
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, id uint) error
}
