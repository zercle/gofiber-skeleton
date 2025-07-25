package repository

import (
	"context"
	"gofiber-skeleton/internal/entities"
	db "gofiber-skeleton/internal/repository/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewURLRepository creates a new URLRepository.
func NewURLRepository(dbpool *pgxpool.Pool) *URLRepository {
	return &URLRepository{queries: db.New(dbpool)}
}

// URLRepository implements the usecases.URLRepository interface.
type URLRepository struct {
	queries *db.Queries
}

// CreateURL creates a new URL in the database.
func (r *URLRepository) CreateURL(ctx context.Context, url *entities.URL) error {
	_, err := r.queries.CreateURL(ctx, db.CreateURLParams{
		OriginalUrl: url.OriginalURL,
		ShortCode:   url.ShortCode,
		UserID:      pgtype.UUID{Bytes: url.UserID, Valid: true},
		ExpiresAt:   pgtype.Timestamptz{Time: url.ExpiresAt, Valid: !url.ExpiresAt.IsZero()},
	})
	return err
}

// GetURLByShortCode retrieves a URL from the database by its short code.
func (r *URLRepository) GetURLByShortCode(ctx context.Context, shortCode string) (*entities.URL, error) {
	url, err := r.queries.GetURLByShortCode(ctx, shortCode)
	if err != nil {
		return nil, err
	}
	return &entities.URL{ID: url.ID.Bytes, OriginalURL: url.OriginalUrl, ShortCode: url.ShortCode, UserID: url.UserID.Bytes, CreatedAt: url.CreatedAt.Time, ExpiresAt: url.ExpiresAt.Time}, nil
}

// GetURLsByUserID retrieves all URLs for a given user ID.
func (r *URLRepository) GetURLsByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.URL, error) {
	urls, err := r.queries.GetURLsByUserID(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return nil, err
	}

	var result []*entities.URL
	for _, url := range urls {
		result = append(result, &entities.URL{ID: url.ID.Bytes, OriginalURL: url.OriginalUrl, ShortCode: url.ShortCode, UserID: url.UserID.Bytes, CreatedAt: url.CreatedAt.Time, ExpiresAt: url.ExpiresAt.Time})
	}
	return result, nil
}

// UpdateURL updates a URL in the database.
func (r *URLRepository) UpdateURL(ctx context.Context, url *entities.URL) error {
	_, err := r.queries.UpdateURL(ctx, db.UpdateURLParams{
		ID:          pgtype.UUID{Bytes: url.ID, Valid: true},
		OriginalUrl: url.OriginalURL,
	})
	return err
}

// DeleteURL deletes a URL from the database.
func (r *URLRepository) DeleteURL(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteURL(ctx, pgtype.UUID{Bytes: id, Valid: true})
}
