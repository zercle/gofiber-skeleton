//go:generate mockgen -source=user_usecase.go -destination=mocks/mock_user_usecase.go -package=mocks UserUseCase

package usecase

import (
	"context"
	"gofiber-skeleton/internal/user"
)

// UserUseCase defines the interface for user-related business logic.
type UserUseCase interface {
	Register(ctx context.Context, username, password string) (*user.ModelUser, error)
	Login(ctx context.Context, username, password string) (string, error) // Returns JWT
}
