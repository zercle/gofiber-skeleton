package usecase

import (
	"context"
	"gofiber-skeleton/internal/user"            // Updated import
	"gofiber-skeleton/internal/user/repository" // Updated import

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// NewUserUseCase creates a new UserUseCase.
import "time"

func NewUserUseCase(userRepo repository.UserRepository, jwtSecret string, jwtExpiration time.Duration) UserUseCase {
	return &userUseCase{userRepo: userRepo, jwtSecret: jwtSecret, jwtExpiration: jwtExpiration}
}

type userUseCase struct {
	userRepo      repository.UserRepository
	jwtSecret     string
	jwtExpiration time.Duration
}

func (uc *userUseCase) Register(ctx context.Context, username, password string) (*user.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	usr := &user.User{ // Changed variable name to avoid conflict with package name
		Username: username,
		Password: string(hashedPassword),
	}

	err = uc.userRepo.CreateUser(ctx, usr) // Changed variable name
	if err != nil {
		return nil, err
	}

	return usr, nil // Changed variable name
}

func (uc *userUseCase) Login(ctx context.Context, username, password string) (string, error) {
	usr, err := uc.userRepo.GetUserByUsername(ctx, username) // Changed variable name
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(password)) // Changed variable name
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
		return "", err
	}

	return t, nil
}
