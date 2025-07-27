package url

import (
	"time"

	"github.com/google/uuid"
)

// URL represents a shortened URL in the system.
type URL struct {
	ID          uuid.UUID `json:"id"`
	OriginalURL string    `json:"original_url"`
	ShortCode   string    `json:"short_code"`
	UserID      uuid.UUID `json:"user_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at,omitempty"`
}
