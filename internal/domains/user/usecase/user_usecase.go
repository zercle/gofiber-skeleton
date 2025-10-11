package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/zercle/gofiber-skeleton/internal/config"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/entity"
	"github.com/zercle/gofiber-skeleton/internal/domains/user/repository"
	"github.com/zercle/gofiber-skeleton/pkg/auth"
	"github.com/zercle/gofiber-skeleton/pkg/validator"
)

// userUsecase implements UserUsecase
type userUsecase struct {
	userRepo   repository.UserRepository
	jwtManager *auth.JWTManager
	validator  *validator.Validator
	cfg        *config.Config
}

// NewUserUsecase creates a new user usecase
func NewUserUsecase(
	userRepo repository.UserRepository,
	jwtManager *auth.JWTManager,
	validator *validator.Validator,
	cfg *config.Config,
) UserUsecase {
	return &userUsecase{
		userRepo:   userRepo,
		jwtManager: jwtManager,
		validator:  validator,
		cfg:        cfg,
	}
}

// Register registers a new user
func (u *userUsecase) Register(ctx context.Context, req *entity.RegisterRequest) (*entity.User, error) {
	// Validate request
	if errs := u.validator.Validate(req); len(errs) > 0 {
		return nil, fmt.Errorf("validation failed: %s", validator.FormatValidationErrors(errs))
	}

	// Check if email already exists
	existingUser, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to check existing email: %w", err)
	}
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Check if username already exists
	existingUser, err = u.userRepo.GetByUsername(ctx, req.Username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to check existing username: %w", err)
	}
	if existingUser != nil {
		return nil, errors.New("username already taken")
	}

	// Hash password
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user entity
	user := &entity.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
		FullName:     req.FullName,
		IsActive:     true,
		IsVerified:   false,
	}

	// Save to database
	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// Login authenticates a user and returns tokens
func (u *userUsecase) Login(ctx context.Context, req *entity.LoginRequest) (*entity.LoginResponse, error) {
	// Validate request
	if errs := u.validator.Validate(req); len(errs) > 0 {
		return nil, fmt.Errorf("validation failed: %s", validator.FormatValidationErrors(errs))
	}

	// Get user by email
	user, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("invalid email or password")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	// Verify password
	if err := auth.CheckPassword(req.Password, user.PasswordHash); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate access token
	accessToken, err := u.jwtManager.GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := u.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Update last login
	if err := u.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		// Log error but don't fail the login
		fmt.Printf("Warning: failed to update last login: %v\n", err)
	}

	return &entity.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    u.jwtManager.GetExpiresIn(),
		User:         user,
	}, nil
}

// GetByID retrieves a user by ID
func (u *userUsecase) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetProfile retrieves the current user's profile
func (u *userUsecase) GetProfile(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	return u.GetByID(ctx, userID)
}

// Update updates a user's profile
func (u *userUsecase) Update(ctx context.Context, userID uuid.UUID, req *entity.UpdateUserRequest) (*entity.User, error) {
	// Validate request
	if errs := u.validator.Validate(req); len(errs) > 0 {
		return nil, fmt.Errorf("validation failed: %s", validator.FormatValidationErrors(errs))
	}

	// Get current user
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Update fields if provided
	if req.FullName != nil {
		user.FullName = req.FullName
	}
	if req.Email != nil && *req.Email != user.Email {
		// Check if new email already exists
		existingUser, err := u.userRepo.GetByEmail(ctx, *req.Email)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to check existing email: %w", err)
		}
		if existingUser != nil {
			return nil, errors.New("email already registered")
		}
		user.Email = *req.Email
	}
	if req.Username != nil && *req.Username != user.Username {
		// Check if new username already exists
		existingUser, err := u.userRepo.GetByUsername(ctx, *req.Username)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to check existing username: %w", err)
		}
		if existingUser != nil {
			return nil, errors.New("username already taken")
		}
		user.Username = *req.Username
	}

	// Save updated user
	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// ChangePassword changes a user's password
func (u *userUsecase) ChangePassword(ctx context.Context, userID uuid.UUID, req *entity.ChangePasswordRequest) error {
	// Validate request
	if errs := u.validator.Validate(req); len(errs) > 0 {
		return fmt.Errorf("validation failed: %s", validator.FormatValidationErrors(errs))
	}

	// Get current user
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("user not found")
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Verify current password
	if err := auth.CheckPassword(req.CurrentPassword, user.PasswordHash); err != nil {
		return errors.New("current password is incorrect")
	}

	// Hash new password
	newPasswordHash, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update password
	if err := u.userRepo.UpdatePassword(ctx, userID, newPasswordHash); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// List retrieves a paginated list of users
func (u *userUsecase) List(ctx context.Context, page, perPage int) (*entity.UserListResponse, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	offset := (page - 1) * perPage

	// Get users
	users, err := u.userRepo.List(ctx, perPage, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	// Get total count
	totalCount, err := u.userRepo.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	return &entity.UserListResponse{
		Users:      users,
		TotalCount: totalCount,
		Page:       page,
		PerPage:    perPage,
	}, nil
}

// Deactivate deactivates a user
func (u *userUsecase) Deactivate(ctx context.Context, userID uuid.UUID) error {
	if err := u.userRepo.Deactivate(ctx, userID); err != nil {
		return fmt.Errorf("failed to deactivate user: %w", err)
	}
	return nil
}

// Activate activates a user
func (u *userUsecase) Activate(ctx context.Context, userID uuid.UUID) error {
	if err := u.userRepo.Activate(ctx, userID); err != nil {
		return fmt.Errorf("failed to activate user: %w", err)
	}
	return nil
}

// Delete permanently deletes a user
func (u *userUsecase) Delete(ctx context.Context, userID uuid.UUID) error {
	if err := u.userRepo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
