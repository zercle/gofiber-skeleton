package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/zercle/gofiber-skeleton/internal/config"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/entity"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/repository"
)

//go:generate mockgen -source=auth.go -destination=mocks/auth.go -package=mocks

// AuthUsecase defines the interface for authentication business logic
type AuthUsecase interface {
	// Register registers a new user
	Register(ctx context.Context, req *entity.RegisterRequest) (*entity.User, error)

	// Login authenticates a user and returns a JWT token
	Login(ctx context.Context, req *entity.LoginRequest) (*entity.LoginResponse, error)

	// GetUserByID retrieves a user by ID
	GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error)

	// GetUserByEmail retrieves a user by email
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)

	// UpdateProfile updates user profile
	UpdateProfile(ctx context.Context, userID uuid.UUID, req *entity.UpdateProfileRequest) (*entity.User, error)

	// ChangePassword changes user password
	ChangePassword(ctx context.Context, userID uuid.UUID, req *entity.ChangePasswordRequest) error

	// VerifyToken verifies a JWT token and returns the user ID
	VerifyToken(tokenString string) (uuid.UUID, error)
}

// authUsecase implements AuthUsecase
type authUsecase struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

// NewAuthUsecase creates a new auth usecase
func NewAuthUsecase(userRepo repository.UserRepository, cfg *config.Config) AuthUsecase {
	return &authUsecase{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

// Register registers a new user
func (u *authUsecase) Register(ctx context.Context, req *entity.RegisterRequest) (*entity.User, error) {
	// Check if user already exists
	existingUser, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Hash password
	hashedPassword, err := u.hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user entity
	user := entity.NewUser(req.Email, hashedPassword, req.FullName)

	// Save to database
	createdUser, err := u.userRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return createdUser, nil
}

// Login authenticates a user and returns a JWT token
func (u *authUsecase) Login(ctx context.Context, req *entity.LoginRequest) (*entity.LoginResponse, error) {
	// Get user by email
	user, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, fmt.Errorf("account is deactivated")
	}

	// Verify password
	if err := u.verifyPassword(user.PasswordHash, req.Password); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token, expiresIn, err := u.generateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &entity.LoginResponse{
		User:        user.ToPublic(),
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int64(expiresIn.Seconds()),
	}, nil
}

// GetUserByID retrieves a user by ID
func (u *authUsecase) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}

// GetUserByEmail retrieves a user by email
func (u *authUsecase) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}

// UpdateProfile updates user profile
func (u *authUsecase) UpdateProfile(ctx context.Context, userID uuid.UUID, req *entity.UpdateProfileRequest) (*entity.User, error) {
	// Get existing user
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Update profile
	user.UpdateProfile(req.FullName)

	// Save changes
	updatedUser, err := u.userRepo.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	return updatedUser, nil
}

// ChangePassword changes user password
func (u *authUsecase) ChangePassword(ctx context.Context, userID uuid.UUID, req *entity.ChangePasswordRequest) error {
	// Get user
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Verify old password
	if err := u.verifyPassword(user.PasswordHash, req.OldPassword); err != nil {
		return fmt.Errorf("invalid old password")
	}

	// Hash new password
	hashedPassword, err := u.hashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	if err := u.userRepo.UpdatePassword(ctx, userID, hashedPassword); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// VerifyToken verifies a JWT token and returns the user ID
func (u *authUsecase) VerifyToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(u.cfg.JWT.Secret), nil
	})

	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			return uuid.Nil, fmt.Errorf("invalid user_id in token")
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return uuid.Nil, fmt.Errorf("invalid user_id format: %w", err)
		}

		return userID, nil
	}

	return uuid.Nil, fmt.Errorf("invalid token claims")
}

// hashPassword hashes a password using bcrypt
func (u *authUsecase) hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// verifyPassword verifies a password against its hash
func (u *authUsecase) verifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// generateToken generates a JWT token for a user
func (u *authUsecase) generateToken(userID uuid.UUID) (string, time.Duration, error) {
	expiresIn := u.cfg.JWT.AccessExpiration
	expirationTime := time.Now().Add(expiresIn)

	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     expirationTime.Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(u.cfg.JWT.Secret))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresIn, nil
}
