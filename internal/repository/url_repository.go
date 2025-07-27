//go:generate mockgen -source=url_repository.go -destination=mocks/mock_redis_cache.go -package=mocks RedisCache

package repository

import (
	"context"
	"gofiber-skeleton/internal/entities"
	db "gofiber-skeleton/internal/repository/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/redis/go-redis/v9"
	"time"
)

// RedisCache defines the interface for Redis operations used by URLRepository.
type RedisCache interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}

// NewURLRepository creates a new URLRepository.
func NewSQLURLRepository(querier DBQueriesInterface, redisClient RedisCache) *URLRepository {
	return &URLRepository{queries: querier, redisClient: redisClient}
}

// URLRepository implements the usecases.URLRepository interface.
type URLRepository struct {
	queries     DBQueriesInterface
	redisClient RedisCache
}

// CreateURL creates a new URL in the database and returns the created URL.
func (r *URLRepository) CreateURL(ctx context.Context, url *entities.URL) (*entities.URL, error) {
	dbURL, err := r.queries.CreateURL(ctx, db.CreateURLParams{
		OriginalUrl: url.OriginalURL,
		ShortCode:   url.ShortCode,
		UserID:      pgtype.UUID{Bytes: url.UserID, Valid: true},
		ExpiresAt:   pgtype.Timestamptz{Time: url.ExpiresAt, Valid: !url.ExpiresAt.IsZero()},
	})
	if err != nil {
		return nil, err
	}

	// Map the database URL to entity URL
	createdURL := &entities.URL{
		ID:          dbURL.ID.Bytes,
		OriginalURL: dbURL.OriginalUrl,
		ShortCode:   dbURL.ShortCode,
		UserID:      dbURL.UserID.Bytes,
		CreatedAt:   dbURL.CreatedAt.Time,
		ExpiresAt:   dbURL.ExpiresAt.Time,
	}

	// Cache the URL
	r.redisClient.Set(ctx, createdURL.ShortCode, createdURL.OriginalURL, 0)

	return createdURL, nil
}

// GetURLByShortCode retrieves a URL from the database by its short code.
func (r *URLRepository) GetURLByShortCode(ctx context.Context, shortCode string) (*entities.URL, error) {
	// Try to get from cache first
	originalURL, err := r.redisClient.Get(ctx, shortCode).Result()
	if err == nil {
		return &entities.URL{OriginalURL: originalURL, ShortCode: shortCode}, nil
	}

	url, err := r.queries.GetURLByShortCode(ctx, shortCode)
	if err != nil {
		return nil, err
	}
	// Cache the URL if found in DB
	r.redisClient.Set(ctx, url.ShortCode, url.OriginalUrl, 0)

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

// UpdateURL updates a URL in the database and returns the updated URL.
func (r *URLRepository) UpdateURL(ctx context.Context, url *entities.URL) (*entities.URL, error) {
	dbURL, err := r.queries.UpdateURL(ctx, db.UpdateURLParams{
		ID:          pgtype.UUID{Bytes: url.ID, Valid: true},
		OriginalUrl: url.OriginalURL,
	})
	if err != nil {
		return nil, err
	}

	// Map the database URL to entity URL
	updatedURL := &entities.URL{
		ID:          dbURL.ID.Bytes,
		OriginalURL: dbURL.OriginalUrl,
		ShortCode:   dbURL.ShortCode, // Keep existing short code if not changed
		UserID:      dbURL.UserID.Bytes,
		CreatedAt:   dbURL.CreatedAt.Time,
		ExpiresAt:   dbURL.ExpiresAt.Time,
	}

	// Update cache
	r.redisClient.Set(ctx, updatedURL.ShortCode, updatedURL.OriginalURL, 0)

	return updatedURL, nil
}

// DeleteURL deletes a URL from the database.
func (r *URLRepository) DeleteURL(ctx context.Context, id uuid.UUID) error {
	// First, get the URL to retrieve its short_code for cache invalidation.
	url, err := r.GetURLByID(ctx, id)
	if err != nil {
		return err // Return error if URL not found or other DB issue.
	}

	// Then, delete the URL from the database.
	if err := r.queries.DeleteURL(ctx, pgtype.UUID{Bytes: id, Valid: true}); err != nil {
		return err
	}

	// Finally, invalidate the cache.
	r.redisClient.Del(ctx, url.ShortCode)

	return nil
}

// GetURLByID retrieves a URL from the database by its ID.
func (r *URLRepository) GetURLByID(ctx context.Context, id uuid.UUID) (*entities.URL, error) {
	url, err := r.queries.GetURLByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}

	return &entities.URL{ID: url.ID.Bytes, OriginalURL: url.OriginalUrl, ShortCode: url.ShortCode, UserID: url.UserID.Bytes, CreatedAt: url.CreatedAt.Time, ExpiresAt: url.ExpiresAt.Time}, nil
}
