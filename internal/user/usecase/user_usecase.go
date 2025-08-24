package userusecase

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/zercle/gofiber-skeleton/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type userUseCase struct {
	userRepo  domain.UserRepository
	jwtSecret string
}

// NewUserUseCase creates a new user use case
func NewUserUseCase(userRepo domain.UserRepository, jwtSecret string) domain.UserUseCase {
	return &userUseCase{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (uc *userUseCase) Register(username, password, role string) (*domain.User, error) {
	// Validate role
	if role != domain.RoleAdmin && role != domain.RoleCustomer {
		role = domain.RoleCustomer // Default to customer
	}

	// Check if user already exists
	existingUser, _ := uc.userRepo.GetByUsername(username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &domain.User{
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

func (uc *userUseCase) Login(username, password string) (string, *domain.User, error) {
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

func (uc *userUseCase) GetByID(id string) (*domain.User, error) {
	return uc.userRepo.GetByID(id)
}

func (uc *userUseCase) UpdateRole(id, role string) error {
	// Validate role
	if role != domain.RoleAdmin && role != domain.RoleCustomer {
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