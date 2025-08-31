package userusecase

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/zercle/gofiber-skeleton/internal/usermodule"
	"golang.org/x/crypto/bcrypt"
)

type userUseCase struct {
	userRepo   usermodule.UserRepository
	jwtSecret  string
	bcryptCost int // New field for bcrypt cost
}

// NewUserUseCase creates a new user use case
func NewUserUseCase(userRepo usermodule.UserRepository, jwtSecret string, bcryptCost int) usermodule.UserUseCase {
	if bcryptCost == 0 { // Default to bcrypt.DefaultCost if not provided
		bcryptCost = bcrypt.DefaultCost
	}
	return &userUseCase{
		userRepo:   userRepo,
		jwtSecret:  jwtSecret,
		bcryptCost: bcryptCost,
	}
}

func (uc *userUseCase) Register(username, password, role string) (*usermodule.User, error) {
	// Validate role
	if role != usermodule.RoleAdmin && role != usermodule.RoleCustomer {
		role = usermodule.RoleCustomer // Default to customer
	}

	// Check if user already exists
	existingUser, _ := uc.userRepo.GetByUsername(username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), uc.bcryptCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &usermodule.User{
		ID:           uuid.New().String(),
		Username:     username,
		PasswordHash: string(hashedPassword),
		Role:         role,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	createdUser, err := uc.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (uc *userUseCase) Login(username, password string) (string, *usermodule.User, error) {
	// Get user by username
	user, err := uc.userRepo.GetByUsername(username)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 24 hours
	})

	tokenString, err := token.SignedString([]byte(uc.jwtSecret))
	if err != nil {
		return "", nil, err
	}

	return tokenString, user, nil
}

func (uc *userUseCase) GetByID(id string) (*usermodule.User, error) {
	return uc.userRepo.GetByID(id)
}

func (uc *userUseCase) UpdateRole(id, role string) error {
	// Validate role
	if role != usermodule.RoleAdmin && role != usermodule.RoleCustomer {
		return errors.New("invalid role")
	}

	user, err := uc.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	user.Role = role
	user.UpdatedAt = time.Now()

	return uc.userRepo.Update(user)
}
