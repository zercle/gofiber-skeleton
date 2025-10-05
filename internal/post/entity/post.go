package entity

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID `json:"id"`
	ThreadID  uuid.UUID `json:"thread_id"`
	UserID    uuid.UUID `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
