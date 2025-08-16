package usecase

import (
	"context"
	"gofiber-skeleton/internal/user"
	"gofiber-skeleton/internal/user/repository"

	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func NewUserUseCase(userRepo repository.UserRepository, jwtSecret string, jwtExpiration time.Duration) UserUseCase {
	return &userUseCase{userRepo: userRepo, jwtSecret: jwtSecret, jwtExpiration: jwtExpiration}
}

type userUseCase struct {
	userRepo      repository.UserRepository
	jwtSecret     string
	jwtExpiration time.Duration
}

func (uc *userUseCase) Register(ctx context.Context, username, password, role string) (*user.ModelUser, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	usr := &user.ModelUser{
		Username: username,
		Password: string(hashedPassword),
		Role:     role,
	}

	err = uc.userRepo.CreateUser(ctx, usr)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (uc *userUseCase) Login(ctx context.Context, username, password string) (string, error) {
	usr, err := uc.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(password))
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"sub": usr.ID,
		"exp": jwt.NewNumericDate(time.Now().Add(uc.jwtExpiration)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(uc.jwtSecret))
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (uc *userUseCase) UpdateRole(ctx context.Context, userID uuid.UUID, role string) (*user.ModelUser, error) {
	usr, err := uc.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	usr.Role = role
	err = uc.userRepo.UpdateUserRole(ctx, userID, role)
	if err != nil {
		return nil, err
	}

	return usr, nil
}
