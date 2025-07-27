package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"gofiber-skeleton/internal/platform/cache"
	"gofiber-skeleton/internal/platform/db"
	"gofiber-skeleton/internal/url"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/redis/go-redis/v9"
)

// NewSQLURLRepository creates a new URLRepository.
func NewSQLURLRepository(queries *db.Queries, redisClient *redis.Client) URLRepository {
	return &SQLURLRepository{
		queries:     queries,
		redisClient: cache.NewCacheWrapper(redisClient),
	}
}

// SQLURLRepository implements the url.URLRepository interface.
type SQLURLRepository struct {
	queries     *db.Queries
	redisClient cache.CacheWrapper
}

// CreateURL creates a new url in the database.
func (r *SQLURLRepository) CreateURL(ctx context.Context, url *url.URL) error {
	_, err := r.queries.CreateURL(ctx, db.CreateURLParams{
		OriginalUrl: url.OriginalURL,
		ShortCode:   url.ShortCode,
		UserID:      pgtype.UUID{Bytes: url.UserID, Valid: url.UserID != uuid.Nil},
		ExpiresAt:   pgtype.Timestamptz{Time: url.ExpiresAt, Valid: !url.ExpiresAt.IsZero()},
	})
	if err != nil {
		return err
	}
	// Cache the URL
	data, err := json.Marshal(url)
	if err == nil {
		if cacheErr := r.redisClient.Set(ctx, r.cacheKeyByShortCode(url.ShortCode), data, time.Hour); cacheErr != nil {
			slog.Warn("Failed to cache URL", "error", cacheErr)
		}
		if cacheErr := r.redisClient.Set(ctx, r.cacheKeyByID(url.ID.String()), data, time.Hour); cacheErr != nil {
			slog.Warn("Failed to cache URL", "error", cacheErr)
		}
	}
	return nil
}

// GetURLByShortCode retrieves a url from the database by its short code with caching.
func (r *SQLURLRepository) GetURLByShortCode(ctx context.Context, shortCode string) (*url.URL, error) {
	cacheKey := r.cacheKeyByShortCode(shortCode)
	cachedData, err := r.redisClient.Get(ctx, cacheKey)
	if err == nil {
		var cachedURL url.URL
		if err := json.Unmarshal([]byte(cachedData), &cachedURL); err == nil {
			return &cachedURL, nil
		}
		// On unmarshal error, fall back to DB
	}

	u, err := r.queries.GetURLByShortCode(ctx, shortCode)
	if err != nil {
		return nil, err
	}
	urlObj := &url.URL{
		ID:          u.ID.Bytes,
		OriginalURL: u.OriginalUrl,
		ShortCode:   u.ShortCode,
		UserID:      u.UserID.Bytes,
		CreatedAt:   u.CreatedAt.Time,
		ExpiresAt:   u.ExpiresAt.Time,
	}
	// Cache the URL
	data, err := json.Marshal(urlObj)
	if err == nil {
		if cacheErr := r.redisClient.Set(ctx, cacheKey, data, time.Hour); cacheErr != nil {
			slog.Warn("Failed to cache URL", "error", cacheErr)
		}
		if cacheErr := r.redisClient.Set(ctx, r.cacheKeyByID(urlObj.ID.String()), data, time.Hour); cacheErr != nil {
			slog.Warn("Failed to cache URL", "error", cacheErr)
		}
	}
	return urlObj, nil
}

// GetURLByID retrieves a url from the database by its ID with caching.
func (r *SQLURLRepository) GetURLByID(ctx context.Context, id uuid.UUID) (*url.URL, error) {
	cacheKey := r.cacheKeyByID(id.String())
	cachedData, err := r.redisClient.Get(ctx, cacheKey)
	if err == nil {
		var cachedURL url.URL
		if err := json.Unmarshal([]byte(cachedData), &cachedURL); err == nil {
			return &cachedURL, nil
		}
		// On unmarshal error, fall back to DB
	}

	u, err := r.queries.GetURLByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}
	urlObj := &url.URL{
		ID:          u.ID.Bytes,
		OriginalURL: u.OriginalUrl,
		ShortCode:   u.ShortCode,
		UserID:      u.UserID.Bytes,
		CreatedAt:   u.CreatedAt.Time,
		ExpiresAt:   u.ExpiresAt.Time,
	}
	// Cache the URL
	data, err := json.Marshal(urlObj)
	if err == nil {
		if cacheErr := r.redisClient.Set(ctx, cacheKey, data, time.Hour); cacheErr != nil {
			slog.Warn("Failed to cache URL", "error", cacheErr)
		}
		if cacheErr := r.redisClient.Set(ctx, r.cacheKeyByShortCode(urlObj.ShortCode), data, time.Hour); cacheErr != nil {
			slog.Warn("Failed to cache URL", "error", cacheErr)
		}
	}
	return urlObj, nil
}

func (r *SQLURLRepository) cacheKeyByShortCode(shortCode string) string {
	return fmt.Sprintf("url:shortcode:%s", shortCode)
}

func (r *SQLURLRepository) cacheKeyByID(id string) string {
	return fmt.Sprintf("url:id:%s", id)
}
