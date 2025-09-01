package usecase

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"github.com/zercle/gofiber-skeleton/internal/user"
)

// userUsecase implements the user.UserUsecase interface.
type userUsecase struct {
	userRepo user.UserRepository
}

// NewUserUsecase creates a new UserUsecase.
func NewUserUsecase(userRepo user.UserRepository) user.UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

// Register creates a new user with a hashed password.
func (uc *userUsecase) Register(payload user.RegisterPayload) (user.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return user.User{}, err
	}

	newUser := user.User{
		Email:     payload.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdUser, err := uc.userRepo.CreateUser(newUser)
	if err != nil {
		return user.User{}, err
	}

	return createdUser, nil
}

func (uc *userUsecase) GetUserByEmail(ctx *fiber.Ctx, email string) (user.User, error) {
	return uc.userRepo.GetUserByEmail(ctx, email)
}