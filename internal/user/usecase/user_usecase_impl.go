package usecase

import (
	"context"
	"errors"
	"gofiber-skeleton/internal/user/domain"
	"github.com/google/uuid"
)

func (uc *userUsecase) Login(ctx context.Context, email, password string) (string, error) {
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		// Return a generic error to avoid user enumeration attacks
		return "", errors.New("invalid credentials")
	}

	// Compare hashed password (e.g., using bcrypt)
	// err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	// if err != nil {
	// 	return "", errors.New("invalid credentials")
	// }

	// On success, generate and return token
	return uc.jwtService.GenerateToken(user.ID.String())
}

func (uu *userUsecase) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return uu.userRepo.GetUser(ctx, id)
}

func (uu *userUsecase) CreateUser(ctx context.Context, user *domain.User) error {
	return uu.userRepo.CreateUser(ctx, user)
}

func (uu *userUsecase) UpdateUser(ctx context.Context, user *domain.User) error {
	return uu.userRepo.UpdateUser(ctx, user)
}

func (uu *userUsecase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return uu.userRepo.DeleteUser(ctx, id)
}
