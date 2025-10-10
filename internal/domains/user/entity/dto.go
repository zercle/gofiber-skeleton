package entity

// RegisterRequest represents user registration request
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name" validate:"required,min=2"`
}

// LoginRequest represents user login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents user login response
type LoginResponse struct {
	User        *PublicUser `json:"user"`
	AccessToken string      `json:"access_token"`
	TokenType   string      `json:"token_type"`
	ExpiresIn   int64       `json:"expires_in"` // seconds
}

// UpdateProfileRequest represents profile update request
type UpdateProfileRequest struct {
	FullName string `json:"full_name" validate:"required,min=2"`
}

// ChangePasswordRequest represents password change request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// UserResponse represents a generic user response
type UserResponse struct {
	User *PublicUser `json:"user"`
}

// UsersResponse represents a list of users response
type UsersResponse struct {
	Users      []*PublicUser `json:"users"`
	Total      int64         `json:"total"`
	Page       int           `json:"page"`
	PerPage    int           `json:"per_page"`
	TotalPages int           `json:"total_pages"`
}
