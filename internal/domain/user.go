package domain

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"unique;not null"`
	Email     string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"` // hashed password
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserRepository interface {
	// Create persists a new user record.
	Create(user *User) error
	// GetByID retrieves a user by its ID.
	GetByID(id uint) (*User, error)
	// GetByUsername retrieves a user by username.
	GetByUsername(username string) (*User, error)
	// GetByEmail retrieves a user by email.
	GetByEmail(email string) (*User, error)
	// Update modifies an existing user.
	Update(user *User) error
	// Delete removes a user by its ID.
	Delete(id uint) error
}