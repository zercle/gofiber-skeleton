package usecase

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"golang.org/x/crypto/argon2"

	"github.com/zercle/gofiber-skeleton/internal/config"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/entity"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/repository"
	"github.com/zercle/gofiber-skeleton/internal/middleware"
)

// UserUsecase defines the interface for user business logic
type UserUsecase interface {
	// Register registers a new user
	Register(ctx context.Context, req *entity.CreateUserRequest) (*entity.UserResponse, error)

	// Login authenticates a user and returns a token
	Login(ctx context.Context, req *entity.LoginRequest) (*entity.LoginResponse, error)

	// GetProfile retrieves user profile by ID
	GetProfile(ctx context.Context, userID uuid.UUID) (*entity.UserResponse, error)

	// UpdateProfile updates user profile
	UpdateProfile(ctx context.Context, userID uuid.UUID, req *entity.UpdateUserRequest) (*entity.UserResponse, error)

	// ChangePassword changes user password
	ChangePassword(ctx context.Context, userID uuid.UUID, req *entity.ChangePasswordRequest) error

	// ListUsers retrieves a list of users with pagination
	ListUsers(ctx context.Context, limit, offset int) ([]*entity.UserResponse, error)

	// DeactivateUser deactivates a user account
	DeactivateUser(ctx context.Context, userID uuid.UUID) error
}

type userUsecase struct {
	userRepo repository.UserRepository
	config   *config.Config
}

// NewUserUsecase creates a new user usecase implementation
func NewUserUsecase(injector do.Injector) (UserUsecase, error) {
	userRepo := do.MustInvoke[repository.UserRepository](injector)
	cfg := do.MustInvoke[*config.Config](injector)

	return &userUsecase{
		userRepo: userRepo,
		config:   cfg,
	}, nil
}

// Register registers a new user
func (u *userUsecase) Register(ctx context.Context, req *entity.CreateUserRequest) (*entity.UserResponse, error) {
	// Validate request
	if err := req.ValidateUserRequest(); err != nil {
		return nil, err
	}

	// Check if email already exists
	exists, err := u.userRepo.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, entity.ErrEmailExists
	}

	// Hash password
	passwordHash, err := u.hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user entity
	user := entity.NewUser(req.Email, passwordHash, req.FullName)

	// Save user
	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user.ToResponse(), nil
}

// Login authenticates a user and returns a token
func (u *userUsecase) Login(ctx context.Context, req *entity.LoginRequest) (*entity.LoginResponse, error) {
	// Validate request
	if err := req.ValidateLoginRequest(); err != nil {
		return nil, err
	}

	// Get user by email
	user, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, entity.ErrInvalidPassword
	}

	// Verify password
	if !u.verifyPassword(user.PasswordHash, req.Password) {
		return nil, entity.ErrInvalidPassword
	}

	// Generate JWT token
	token, err := u.generateJWT(user)
	if err != nil {
		return nil, err
	}

	return &entity.LoginResponse{
		Token: token,
		User:  user.ToResponse(),
	}, nil
}

// GetProfile retrieves user profile by ID
func (u *userUsecase) GetProfile(ctx context.Context, userID uuid.UUID) (*entity.UserResponse, error) {
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user.ToResponse(), nil
}

// UpdateProfile updates user profile
func (u *userUsecase) UpdateProfile(ctx context.Context, userID uuid.UUID, req *entity.UpdateUserRequest) (*entity.UserResponse, error) {
	// Get existing user
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Email != "" && req.Email != user.Email {
		// Check if new email already exists
		exists, err := u.userRepo.EmailExists(ctx, req.Email)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, entity.ErrEmailExists
		}
		user.Email = req.Email
	}

	if req.FullName != "" {
		user.FullName = req.FullName
	}

	user.UpdatedAt = time.Now()

	// Update user
	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user.ToResponse(), nil
}

// ChangePassword changes user password
func (u *userUsecase) ChangePassword(ctx context.Context, userID uuid.UUID, req *entity.ChangePasswordRequest) error {
	// Get user
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify current password
	if !u.verifyPassword(user.PasswordHash, req.CurrentPassword) {
		return entity.ErrInvalidPassword
	}

	// Hash new password
	newPasswordHash, err := u.hashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// Update password
	return u.userRepo.UpdatePassword(ctx, userID, newPasswordHash)
}

// ListUsers retrieves a list of users with pagination
func (u *userUsecase) ListUsers(ctx context.Context, limit, offset int) ([]*entity.UserResponse, error) {
	users, err := u.userRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	responses := make([]*entity.UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}

	return responses, nil
}

// DeactivateUser deactivates a user account
func (u *userUsecase) DeactivateUser(ctx context.Context, userID uuid.UUID) error {
	return u.userRepo.Deactivate(ctx, userID)
}

// Password hashing and verification methods

func (u *userUsecase) hashPassword(password string) (string, error) {
	// Generate a random salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// Hash the password using Argon2id
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	// Return the salt and hash combined
	return base64.RawStdEncoding.EncodeToString(append(salt, hash...)), nil
}

func (u *userUsecase) verifyPassword(hashedPassword, password string) bool {
	// Decode the stored password
	data, err := base64.RawStdEncoding.DecodeString(hashedPassword)
	if err != nil || len(data) < 16 {
		return false
	}

	// Extract salt and hash
	salt := data[:16]
	hash := data[16:]

	// Hash the provided password with the same salt
	computedHash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	// Compare the hashes
	return subtle.ConstantTimeCompare(hash, computedHash) == 1
}

// JWT generation method

func (u *userUsecase) generateJWT(user *entity.DomainUser) (string, error) {
	// Parse JWT expiry duration
	expiry, err := time.ParseDuration(u.config.JWT.Expiry)
	if err != nil {
		expiry = 24 * time.Hour // default to 24 hours
	}

	// Create claims
	claims := &middleware.Claims{
		UserID: user.ID.String(),
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    u.config.AppName,
			Subject:   user.ID.String(),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	return token.SignedString([]byte(u.config.JWT.Secret))
}
