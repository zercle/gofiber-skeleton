package models

import (
	"time"

	"github.com/google/uuid"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         UserData  `json:"user"`
}

type RegisterRequest struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
	FirstName       string `json:"first_name" validate:"required,min=2,max=50"`
	LastName        string `json:"last_name" validate:"required,min=2,max=50"`
}

type RegisterResponse struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         UserData  `json:"user"`
}

type UserData struct {
	ID              uuid.UUID  `json:"id"`
	Email           string     `json:"email"`
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	FullName        string     `json:"full_name"`
	IsActive        bool       `json:"is_active"`
	IsEmailVerified bool       `json:"is_email_verified"`
	LastLoginAt     *time.Time `json:"last_login_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RefreshTokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}
