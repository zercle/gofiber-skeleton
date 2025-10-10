package entity

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Never expose password hash in JSON
	FullName     string    `json:"full_name"`
	IsActive     bool      `json:"is_active"`
	IsVerified   bool      `json:"is_verified"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// NewUser creates a new user instance
func NewUser(email, passwordHash, fullName string) *User {
	return &User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: passwordHash,
		FullName:     fullName,
		IsActive:     true,
		IsVerified:   false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// IsValidForRegistration checks if user data is valid for registration
func (u *User) IsValidForRegistration() bool {
	return u.Email != "" && u.PasswordHash != "" && u.FullName != ""
}

// Activate activates the user account
func (u *User) Activate() {
	u.IsActive = true
	u.UpdatedAt = time.Now()
}

// Deactivate deactivates the user account
func (u *User) Deactivate() {
	u.IsActive = false
	u.UpdatedAt = time.Now()
}

// Verify marks the user as verified
func (u *User) Verify() {
	u.IsVerified = true
	u.UpdatedAt = time.Now()
}

// UpdatePassword updates the user's password hash
func (u *User) UpdatePassword(newPasswordHash string) {
	u.PasswordHash = newPasswordHash
	u.UpdatedAt = time.Now()
}

// UpdateProfile updates the user's profile information
func (u *User) UpdateProfile(fullName string) {
	if fullName != "" {
		u.FullName = fullName
	}
	u.UpdatedAt = time.Now()
}

// PublicUser returns a user struct without sensitive information
type PublicUser struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	FullName   string    `json:"full_name"`
	IsActive   bool      `json:"is_active"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
}

// ToPublic converts User to PublicUser
func (u *User) ToPublic() *PublicUser {
	return &PublicUser{
		ID:         u.ID,
		Email:      u.Email,
		FullName:   u.FullName,
		IsActive:   u.IsActive,
		IsVerified: u.IsVerified,
		CreatedAt:  u.CreatedAt,
	}
}
