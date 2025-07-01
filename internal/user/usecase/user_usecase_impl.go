package usecase

import (
	"context"
	"gofiber-skeleton/internal/user/domain"
	"gofiber-skeleton/internal/user/infrastructure"
)

type userUsecase struct {
	userRepo infrastructure.UserRepository
}

func NewUserUsecase(userRepo infrastructure.UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

func (uu *userUsecase) GetUser(ctx context.Context, id uint) (*domain.User, error) {
	return uu.userRepo.GetUser(ctx, id)
}

func (uu *userUsecase) CreateUser(ctx context.Context, user *domain.User) error {
	return uu.userRepo.CreateUser(ctx, user)
}

func (uu *userUsecase) UpdateUser(ctx context.Context, user *domain.User) error {
	return uu.userRepo.UpdateUser(ctx, user)
}

func (uu *userUsecase) DeleteUser(ctx context.Context, id uint) error {
	return uu.userRepo.DeleteUser(ctx, id)
}
