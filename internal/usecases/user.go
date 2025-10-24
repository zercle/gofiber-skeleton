package usecases

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/zercle/template-go-fiber/internal/domains"
	"github.com/zercle/template-go-fiber/internal/errors"
)

//go:generate mockgen -source=user.go -destination=../../test/unit/mocks/mock_user_usecase.go -package=mocks

// UserUsecase implements the user business logic
type UserUsecase struct {
	repo domains.UserRepository
}

// NewUserUsecase creates a new user usecase
func NewUserUsecase(repo domains.UserRepository) domains.UserUsecase {
	return &UserUsecase{
		repo: repo,
	}
}

// GetUserByID retrieves a user by ID
func (u *UserUsecase) GetUserByID(ctx context.Context, id string) (*domains.User, error) {
	user, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.NewDatabaseError("failed to fetch user", err)
	}

	if user == nil {
		return nil, errors.NewNotFoundError(fmt.Sprintf("user with id %s not found", id))
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email
func (u *UserUsecase) GetUserByEmail(ctx context.Context, email string) (*domains.User, error) {
	user, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.NewDatabaseError("failed to fetch user", err)
	}

	if user == nil {
		return nil, errors.NewNotFoundError(fmt.Sprintf("user with email %s not found", email))
	}

	return user, nil
}

// RegisterUser creates a new user account
func (u *UserUsecase) RegisterUser(ctx context.Context, input *domains.RegisterUserInput) (*domains.User, error) {
	// Validate input
	if input == nil {
		return nil, errors.NewValidationError("registration input is required", nil)
	}

	if input.Email == "" {
		return nil, errors.NewValidationError("email is required", nil)
	}

	if input.Password == "" {
		return nil, errors.NewValidationError("password is required", nil)
	}

	// Check if email already exists
	existingUser, err := u.repo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, errors.NewDatabaseError("failed to check email", err)
	}

	if existingUser != nil {
		return nil, errors.NewDuplicateEntryError(fmt.Sprintf("email %s is already registered", input.Email))
	}

	// Create new user
	user := &domains.User{
		ID:           uuid.New().String(),
		Email:        input.Email,
		PasswordHash: hashPassword(input.Password), // In real app, use proper password hashing
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		IsActive:     true,
	}

	if err := u.repo.Create(ctx, user); err != nil {
		return nil, errors.NewDatabaseError("failed to create user", err)
	}

	return user, nil
}

// UpdateUser updates user information
func (u *UserUsecase) UpdateUser(ctx context.Context, id string, input *domains.UpdateUserInput) (*domains.User, error) {
	// Validate input
	if input == nil {
		return nil, errors.NewValidationError("update input is required", nil)
	}

	// Get existing user
	user, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.NewDatabaseError("failed to fetch user", err)
	}

	if user == nil {
		return nil, errors.NewNotFoundError(fmt.Sprintf("user with id %s not found", id))
	}

	// Update fields
	if input.Email != nil {
		// Check if new email is already taken
		existingUser, err := u.repo.GetByEmail(ctx, *input.Email)
		if err != nil {
			return nil, errors.NewDatabaseError("failed to check email", err)
		}

		if existingUser != nil && existingUser.ID != user.ID {
			return nil, errors.NewDuplicateEntryError(fmt.Sprintf("email %s is already in use", *input.Email))
		}

		user.Email = *input.Email
	}

	if input.FirstName != nil {
		user.FirstName = input.FirstName
	}

	if input.LastName != nil {
		user.LastName = input.LastName
	}

	if input.IsActive != nil {
		user.IsActive = *input.IsActive
	}

	if err := u.repo.Update(ctx, user); err != nil {
		return nil, errors.NewDatabaseError("failed to update user", err)
	}

	return user, nil
}

// DeleteUser deletes a user account
func (u *UserUsecase) DeleteUser(ctx context.Context, id string) error {
	// Check if user exists
	user, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return errors.NewDatabaseError("failed to fetch user", err)
	}

	if user == nil {
		return errors.NewNotFoundError(fmt.Sprintf("user with id %s not found", id))
	}

	if err := u.repo.Delete(ctx, id); err != nil {
		return errors.NewDatabaseError("failed to delete user", err)
	}

	return nil
}

// ListUsers retrieves paginated user list
func (u *UserUsecase) ListUsers(ctx context.Context, limit, offset int32) ([]*domains.User, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	if offset < 0 {
		offset = 0
	}

	users, err := u.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, errors.NewDatabaseError("failed to fetch users", err)
	}

	return users, nil
}

// Helper function to hash password (placeholder - use proper hashing in production)
func hashPassword(password string) string {
	// In production, use bcrypt or argon2
	return fmt.Sprintf("hashed_%s", password)
}
