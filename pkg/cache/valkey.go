package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zercle/gofiber-skeleton/internal/config"
)

// Valkey represents a Valkey/Redis cache client
type Valkey struct {
	Client *redis.Client
}

// NewValkey creates a new Valkey/Redis client
func NewValkey(cfg *config.Config) (*Valkey, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.GetCacheAddress(),
		Password:     cfg.Cache.Password,
		DB:           cfg.Cache.DB,
		PoolSize:     cfg.Cache.PoolSize,
		MinIdleConns: 5,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  4 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Verify connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to cache: %w", err)
	}

	return &Valkey{Client: client}, nil
}

// Close closes the cache client connection
func (v *Valkey) Close() error {
	if v.Client != nil {
		return v.Client.Close()
	}
	return nil
}

// Ping checks if the cache connection is alive
func (v *Valkey) Ping(ctx context.Context) error {
	return v.Client.Ping(ctx).Err()
}

// Set stores a key-value pair with expiration
func (v *Valkey) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return v.Client.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a value by key
func (v *Valkey) Get(ctx context.Context, key string) (string, error) {
	return v.Client.Get(ctx, key).Result()
}

// Delete removes a key
func (v *Valkey) Delete(ctx context.Context, keys ...string) error {
	return v.Client.Del(ctx, keys...).Err()
}

// Exists checks if a key exists
func (v *Valkey) Exists(ctx context.Context, keys ...string) (int64, error) {
	return v.Client.Exists(ctx, keys...).Result()
}

// Expire sets expiration on a key
func (v *Valkey) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return v.Client.Expire(ctx, key, expiration).Err()
}

// Health returns the health status of the cache
func (v *Valkey) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := v.Ping(ctx); err != nil {
		return fmt.Errorf("cache health check failed: %w", err)
	}

	return nil
}
