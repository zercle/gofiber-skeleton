package user

import (
	"time"

	"github.com/google/uuid"
)

// ModelUser represents a user in the system.
type ModelUser struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
