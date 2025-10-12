package entity

import (
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/zercle/gofiber-skeleton/pkg/uuid"
	uuidv7 "github.com/google/uuid"
)

// DomainUser represents a user entity in the domain
type DomainUser struct {
	ID           uuidv7.UUID `json:"id" db:"id"`
	Email        string      `json:"email" db:"email"`
	PasswordHash string      `json:"-" db:"password_hash"` // Hidden in JSON
	FullName     string      `json:"full_name" db:"full_name"`
	IsActive     bool        `json:"is_active" db:"is_active"`
	CreatedAt    time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at" db:"updated_at"`
}

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=128"`
	FullName string `json:"full_name" validate:"required,min=2,max=255"`
}

// LoginRequest represents the login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=1,max=128"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token string        `json:"token"`
	User  *UserResponse `json:"user"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	Email    string `json:"email,omitempty" validate:"omitempty,email,max=255"`
	FullName string `json:"full_name,omitempty" validate:"omitempty,min=2,max=255"`
}

// ChangePasswordRequest represents the request to change password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required,min=1,max=128"`
	NewPassword     string `json:"new_password" validate:"required,min=8,max=128"`
}

// UserResponse represents the user response (without sensitive data)
type UserResponse struct {
	ID        uuidv7.UUID `json:"id"`
	Email     string      `json:"email"`
	FullName  string      `json:"full_name"`
	IsActive  bool        `json:"is_active"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// NewUser creates a new user entity
func NewUser(email, passwordHash, fullName string) *DomainUser {
	return &DomainUser{
		ID:           uuid.NewV7(),
		Email:        email,
		PasswordHash: passwordHash,
		FullName:     fullName,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// Convert between User (sqlc generated) and DomainUser
func UserToDomainUser(sqlcUser User) *DomainUser {
	return &DomainUser{
		ID:           uuidv7.UUID(sqlcUser.ID.Bytes),
		Email:        sqlcUser.Email,
		PasswordHash: sqlcUser.PasswordHash,
		FullName:     sqlcUser.FullName,
		IsActive:     sqlcUser.IsActive.Bool,
		CreatedAt:    sqlcUser.CreatedAt.Time,
		UpdatedAt:    sqlcUser.UpdatedAt.Time,
	}
}

// DomainUserToUser converts DomainUser to User (sqlc type)
func DomainUserToUser(u *DomainUser) User {
	var userUUID [16]byte
	copy(userUUID[:], u.ID[:])
	return User{
		ID:           pgtype.UUID{Bytes: userUUID, Valid: true},
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		FullName:     u.FullName,
		IsActive:     pgtype.Bool{Bool: u.IsActive, Valid: true},
		CreatedAt:    pgtype.Timestamptz{Time: u.CreatedAt, Valid: true},
		UpdatedAt:    pgtype.Timestamptz{Time: u.UpdatedAt, Valid: true},
	}
}

// ToResponse converts DomainUser to UserResponse (hides sensitive data)
func (u *DomainUser) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		FullName:  u.FullName,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// ValidateUserRequest validates user creation request
func (r *CreateUserRequest) ValidateUserRequest() error {
	// Basic validation would be handled by validator middleware
	// Additional business logic validation can be added here
	if r.Email == "" {
		return ErrEmailRequired
	}
	if r.Password == "" {
		return ErrPasswordRequired
	}
	if r.FullName == "" {
		return ErrFullNameRequired
	}
	return nil
}

// ValidateLoginRequest validates login request
func (r *LoginRequest) ValidateLoginRequest() error {
	if r.Email == "" {
		return ErrEmailRequired
	}
	if r.Password == "" {
		return ErrPasswordRequired
	}
	return nil
}

// Error definitions
var (
	ErrEmailRequired    = errors.New("email is required")
	ErrPasswordRequired = errors.New("password is required")
	ErrFullNameRequired = errors.New("full name is required")
	ErrInvalidEmail     = errors.New("invalid email format")
	ErrUserNotFound     = errors.New("user not found")
	ErrUserInactive     = errors.New("user account is inactive")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrEmailExists      = errors.New("email already exists")
	ErrWeakPassword     = errors.New("password is too weak")
)

// DomainError represents a domain-specific error
type DomainError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// NewDomainError creates a new domain error
func NewDomainError(code, message string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
	}
}

// Error implements the error interface
func (e *DomainError) Error() string {
	return e.Message
}
