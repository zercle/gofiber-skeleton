package entity

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user entity in the domain
type User struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	Email        string     `json:"email" db:"email"`
	PasswordHash string     `json:"-" db:"password_hash"` // Hidden in JSON
	FullName     string     `json:"full_name" db:"full_name"`
	IsActive     bool       `json:"is_active" db:"is_active"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
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
	Token string       `json:"token"`
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
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewUser creates a new user entity
func NewUser(email, passwordHash, fullName string) *User {
	return &User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: passwordHash,
		FullName:     fullName,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// ToResponse converts User to UserResponse (hides sensitive data)
func (u *User) ToResponse() *UserResponse {
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
	ErrEmailRequired    = NewDomainError("EMAIL_REQUIRED", "Email is required")
	ErrPasswordRequired = NewDomainError("PASSWORD_REQUIRED", "Password is required")
	ErrFullNameRequired = NewDomainError("FULL_NAME_REQUIRED", "Full name is required")
	ErrInvalidEmail     = NewDomainError("INVALID_EMAIL", "Invalid email format")
	ErrUserNotFound     = NewDomainError("USER_NOT_FOUND", "User not found")
	ErrUserInactive     = NewDomainError("USER_INACTIVE", "User account is inactive")
	ErrInvalidPassword  = NewDomainError("INVALID_PASSWORD", "Invalid password")
	ErrEmailExists      = NewDomainError("EMAIL_EXISTS", "Email already exists")
	ErrWeakPassword     = NewDomainError("WEAK_PASSWORD", "Password is too weak")
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