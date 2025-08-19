package usecase

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	"github.com/zercle/gofiber-skeleton/internal/domain"
)

//go:generate mockgen -source=user_usecase.go -destination=../domain/mock_user_usecase.go -package=domain

type UserUsecase interface {
	CreateUser(ctx context.Context, username, password, role string) (domain.User, error)
	Login(ctx context.Context, username, password string) (domain.User, error)
}

type userUsecase struct {
	repo domain.UserRepository
}

func NewUserUsecase(repo domain.UserRepository) UserUsecase {
	return &userUsecase{repo: repo}
}

func (u *userUsecase) CreateUser(ctx context.Context, username, password, role string) (domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, fiber.NewError(fiber.StatusInternalServerError, "failed to hash password")
	}

	user := domain.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
		Role:         role,
	}

	createdUser, err := u.repo.CreateUser(ctx, user)
	if err != nil {
		return domain.User{}, fiber.NewError(fiber.StatusInternalServerError, "failed to create user")
	}
	return createdUser, nil
}

func (u *userUsecase) Login(ctx context.Context, username, password string) (domain.User, error) {
	user, err := u.repo.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.User{}, fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
		}
		return domain.User{}, fiber.NewError(fiber.StatusInternalServerError, "failed to retrieve user")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return domain.User{}, fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
	}

	return user, nil
}