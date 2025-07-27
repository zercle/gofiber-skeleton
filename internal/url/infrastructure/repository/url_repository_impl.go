package repository

import (
	"context"
	"gofiber-skeleton/internal/platform/db"
	"gofiber-skeleton/internal/url"
	"gofiber-skeleton/internal/url/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// NewSQLURLRepository creates a new URLRepository.
func NewSQLURLRepository(queries *db.Queries) repository.URLRepository {
	return &SQLURLRepository{queries: queries}
}

// SQLURLRepository implements the url.URLRepository interface.
type SQLURLRepository struct {
	queries *db.Queries
}

// CreateURL creates a new url in the database.
func (r *SQLURLRepository) CreateURL(ctx context.Context, url *url.URL) error {
	_, err := r.queries.CreateURL(ctx, db.CreateURLParams{
		OriginalUrl: url.OriginalURL,
		ShortCode:   url.ShortCode,
		UserID:      pgtype.UUID{Bytes: url.UserID, Valid: url.UserID != uuid.Nil},
		ExpiresAt:   pgtype.Timestamptz{Time: url.ExpiresAt, Valid: !url.ExpiresAt.IsZero()},
	})
	return err
}

// GetURLByShortCode retrieves a url from the database by its short code.
func (r *SQLURLRepository) GetURLByShortCode(ctx context.Context, shortCode string) (*url.URL, error) {
	u, err := r.queries.GetURLByShortCode(ctx, shortCode)
	if err != nil {
		return nil, err
	}

	return &url.URL{
		ID:          u.ID.Bytes,
		OriginalURL: u.OriginalUrl,
		ShortCode:   u.ShortCode,
		UserID:      u.UserID.Bytes,
		CreatedAt:   u.CreatedAt.Time,
		ExpiresAt:   u.ExpiresAt.Time,
	}, nil
}

// GetURLByID retrieves a url from the database by its ID.
func (r *SQLURLRepository) GetURLByID(ctx context.Context, id uuid.UUID) (*url.URL, error) {
	u, err := r.queries.GetURLByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}

	return &url.URL{
		ID:          u.ID.Bytes,
		OriginalURL: u.OriginalUrl,
		ShortCode:   u.ShortCode,
		UserID:      u.UserID.Bytes,
		CreatedAt:   u.CreatedAt.Time,
		ExpiresAt:   u.ExpiresAt.Time,
	}, nil
}
