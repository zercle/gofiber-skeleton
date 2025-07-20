package entities

import "time"

type URL struct {
	ID          string    `json:"id"`
	UserID      *string   `json:"user_id,omitempty"`
	ShortCode   string    `json:"short_code"`
	OriginalURL string    `json:"original_url"`
	CreatedAt   time.Time `json:"created_at"`
}
