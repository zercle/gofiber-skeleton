package domain

import (
	"time"

	"github.com/google/uuid"
)

// UserRole represents the role of a user in the system
type UserRole string

const (
	UserRoleAdmin    UserRole = "admin"
	UserRoleCustomer UserRole = "customer"
)

// User represents a user in the e-commerce system
type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // Don't expose password in JSON
	Role      UserRole  `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	Username string   `json:"username" validate:"required,min=3,max=50"`
	Password string   `json:"password" validate:"required,min=6"`
	Role     UserRole `json:"role" validate:"required,oneof=admin customer"`
}

// LoginRequest represents the request to authenticate a user
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents the response after successful authentication
type LoginResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

// UpdateUserRoleRequest represents the request to update a user's role
type UpdateUserRoleRequest struct {
	Role UserRole `json:"role" validate:"required,oneof=admin customer"`
}

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(user *User) error
	GetByID(id uuid.UUID) (*User, error)
	GetByUsername(username string) (*User, error)
	UpdateRole(id uuid.UUID, role UserRole) error
}

// UserService defines the interface for user business logic
type UserService interface {
	Register(req *CreateUserRequest) (*User, error)
	Login(req *LoginRequest) (*LoginResponse, error)
	GetUser(id uuid.UUID) (*User, error)
	UpdateUserRole(id uuid.UUID, req *UpdateUserRoleRequest) (*User, error)
}