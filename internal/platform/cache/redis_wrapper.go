//go:generate mockgen -source=redis_wrapper.go -destination=mocks/mock_redis_wrapper.go -package=mocks CacheWrapper
package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var _ CacheWrapper = (*cacheWrapper)(nil)

// CacheWrapper interface defines caching methods needed.
type CacheWrapper interface {
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, keys ...string) error
}

type cacheWrapper struct {
	Client *redis.Client
}

func NewCacheWrapper(client *redis.Client) CacheWrapper {
	return &cacheWrapper{Client: client}
}

func (r *cacheWrapper) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

func (r *cacheWrapper) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

func (r *cacheWrapper) Del(ctx context.Context, keys ...string) error {
	return r.Client.Del(ctx, keys...).Err()
}
