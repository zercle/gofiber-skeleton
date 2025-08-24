//go:generate mockgen -source=user.go -destination=./mock/mock_user.go -package=mock
package domain

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"` // Never expose password hash in JSON
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// UserRole constants
const (
	RoleAdmin    = "admin"
	RoleCustomer = "customer"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(user *User) (*User, error)
	GetByID(id string) (*User, error)
	GetByUsername(username string) (*User, error)
	Update(user *User) error
	Delete(id string) error
}

// UserUseCase defines the interface for user business logic
type UserUseCase interface {
	Register(username, password, role string) (*User, error)
	Login(username, password string) (string, *User, error)
	GetByID(id string) (*User, error)
	UpdateRole(id, role string) error
}
