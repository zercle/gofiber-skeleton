package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/zercle/gofiber-skeleton/pkg/domains/auth/entities"
	"github.com/zercle/gofiber-skeleton/pkg/domains/auth/models"
	"github.com/zercle/gofiber-skeleton/pkg/domains/auth/repositories"
	"github.com/zercle/gofiber-skeleton/internal/infrastructure/config"
)

type AuthUseCase interface {
	Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error)
	Register(ctx context.Context, req *models.RegisterRequest) (*models.RegisterResponse, error)
	RefreshToken(ctx context.Context, req *models.RefreshTokenRequest) (*models.RefreshTokenResponse, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, req *models.ChangePasswordRequest) error
	GetProfile(ctx context.Context, userID uuid.UUID) (*models.UserData, error)
}

type authUseCase struct {
	userRepo repositories.UserRepository
	cfg      *config.Config
}

// NewAuthUseCase creates a new AuthUseCase.
func NewAuthUseCase(userRepo repositories.UserRepository, cfg *config.Config) AuthUseCase {
	return &authUseCase{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (a *authUseCase) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	user, err := a.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	// Update last login
	if err := a.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		// Log error but don't fail the login
	}

	token, expiresAt, err := a.generateToken(user)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	refreshToken, err := a.generateRefreshToken(user)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	return &models.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User:         a.toUserData(user),
	}, nil
}

func (a *authUseCase) Register(ctx context.Context, req *models.RegisterRequest) (*models.RegisterResponse, error) {
	// Check if user already exists
	existingUser, _ := a.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user
	user := &entities.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		IsActive:     true,
	}

	createdUser, err := a.userRepo.Create(ctx, user)
	if err != nil {
		return nil, errors.New("failed to create user")
	}

	token, expiresAt, err := a.generateToken(createdUser)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	refreshToken, err := a.generateRefreshToken(createdUser)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	return &models.RegisterResponse{
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User:         a.toUserData(createdUser),
	}, nil
}

func (a *authUseCase) RefreshToken(ctx context.Context, req *models.RefreshTokenRequest) (*models.RefreshTokenResponse, error) {
	// Parse refresh token
	token, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.cfg.JWT.Secret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("invalid user ID in token")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	// Get user to generate new token
	user, err := a.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	newToken, expiresAt, err := a.generateToken(user)
	if err != nil {
		return nil, errors.New("failed to generate new token")
	}

	return &models.RefreshTokenResponse{
		Token:     newToken,
		ExpiresAt: expiresAt,
	}, nil
}

func (a *authUseCase) ChangePassword(ctx context.Context, userID uuid.UUID, req *models.ChangePasswordRequest) error {
	user, err := a.userRepo.GetByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash new password")
	}

	user.PasswordHash = string(hashedPassword)
	_, err = a.userRepo.Update(ctx, user)
	if err != nil {
		return errors.New("failed to update password")
	}

	return nil
}

func (a *authUseCase) GetProfile(ctx context.Context, userID uuid.UUID) (*models.UserData, error) {
	user, err := a.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &models.UserData{
		ID:              user.ID,
		Email:           user.Email,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		FullName:        user.FullName(),
		IsActive:        user.IsActive,
		IsEmailVerified: user.IsEmailVerified,
		LastLoginAt:     user.LastLoginAt,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
	}, nil
}

func (a *authUseCase) generateToken(user *entities.User) (string, time.Time, error) {
	expirationTime, err := time.ParseDuration(a.cfg.JWT.ExpiresIn)
	if err != nil {
		expirationTime = 24 * time.Hour // Default to 24 hours
	}

	expiresAt := time.Now().Add(expirationTime)

	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"exp":     expiresAt.Unix(),
		"iat":     time.Now().Unix(),
		"iss":     a.cfg.JWT.Issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(a.cfg.JWT.Secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

func (a *authUseCase) generateRefreshToken(user *entities.User) (string, error) {
	// Refresh tokens have longer expiration
	expiresAt := time.Now().Add(7 * 24 * time.Hour) // 7 days

	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"type":    "refresh",
		"exp":     expiresAt.Unix(),
		"iat":     time.Now().Unix(),
		"iss":     a.cfg.JWT.Issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.cfg.JWT.Secret))
}

func (a *authUseCase) toUserData(user *entities.User) models.UserData {
	return models.UserData{
		ID:              user.ID,
		Email:           user.Email,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		FullName:        user.FullName(),
		IsActive:        user.IsActive,
		IsEmailVerified: user.IsEmailVerified,
		LastLoginAt:     user.LastLoginAt,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
	}
}
