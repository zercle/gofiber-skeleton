package usecase

import (
	"context"
	"gofiber-skeleton/internal/user/domain"
	"gofiber-skeleton/internal/user/repository"
	"gofiber-skeleton/internal/infra/types"
)

type userUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

func (uu *userUsecase) GetUser(ctx context.Context, id types.UUIDv7) (*domain.User, error) {
	return uu.userRepo.GetUser(ctx, id)
}

func (uu *userUsecase) CreateUser(ctx context.Context, user *domain.User) error {
	return uu.userRepo.CreateUser(ctx, user)
}

func (uu *userUsecase) UpdateUser(ctx context.Context, user *domain.User) error {
	return uu.userRepo.UpdateUser(ctx, user)
}

func (uu *userUsecase) DeleteUser(ctx context.Context, id types.UUIDv7) error {
	return uu.userRepo.DeleteUser(ctx, id)
}
