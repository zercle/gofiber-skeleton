//go:generate mockgen -source=user_usecase.go -destination=../mocks/mock_user_usecase.go -package=mocks

package usecase

import (
	"context"
	"gofiber-skeleton/internal/infra/auth"
	"gofiber-skeleton/internal/user/domain"
	"gofiber-skeleton/internal/user/repository"

	"github.com/google/uuid"
)

type UserUsecase interface {
	Login(ctx context.Context, email, password string) (string, error)
	GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) error
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type userUsecase struct {
	userRepo   repository.UserRepository
	jwtService auth.JWTService
}

func NewUserUsecase(userRepo repository.UserRepository, jwtService auth.JWTService) UserUsecase {
	return &userUsecase{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}
