package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/zercle/gofiber-skeleton/internal/user/entity"
	"github.com/zercle/gofiber-skeleton/internal/user/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Register(ctx context.Context, username, email, password string) (*entity.User, error)
	Login(ctx context.Context, email, password string) (token string, user *entity.User, err error)
}

type authUsecase struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

func NewAuthUsecase(userRepo repository.UserRepository, jwtSecret string) AuthUsecase {
	return &authUsecase{userRepo: userRepo, jwtSecret: jwtSecret}
}

func (uc *authUsecase) Register(ctx context.Context, username, email, password string) (*entity.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user := &entity.User{
		Username:  username,
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if user.ID, err = uuid.NewV7(); err != nil {
		return nil, err
	}

	err = uc.userRepo.Create(ctx, user, string(hashedPassword))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *authUsecase) Login(ctx context.Context, email, password string) (token string, user *entity.User, err error) {
	user, hashedPassword, err := uc.userRepo.GetByEmail(ctx, email)
	if err == sql.ErrNoRows {
		return "", nil, fmt.Errorf("user with email %s not found", email)
	}
	if err != nil {
		return "", nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return "", nil, err
	}

	// Create JWT token
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = jwtToken.SignedString([]byte(uc.jwtSecret))
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}
