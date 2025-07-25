package usecases

import (
	"context"
	"gofiber-skeleton/internal/entities"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// NewUserUseCase creates a new UserUseCase.
func NewUserUseCase(userRepo UserRepository, jwtSecret string, jwtExpiration int) UserUseCase {
	return &userUseCase{userRepo: userRepo, jwtSecret: jwtSecret, jwtExpiration: jwtExpiration}
}

type userUseCase struct {
	userRepo      UserRepository
	jwtSecret     string
	jwtExpiration int
}

func (uc *userUseCase) Register(ctx context.Context, username, password string) (*entities.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &entities.User{
		Username: username,
		Password: string(hashedPassword),
	}

	err = uc.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *userUseCase) Login(ctx context.Context, username, password string) (string, error) {
	user, err := uc.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(uc.jwtExpiration))),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(uc.jwtSecret))
	if err != nil {
		return "", err
	}

	return t, nil
}
