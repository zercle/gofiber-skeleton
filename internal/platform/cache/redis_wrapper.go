package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RealRedisClient implements the RedisCache interface using a real redis.Client.
type RealRedisClient struct {
	Client *redis.Client
}

// Set implements the Set method of the RedisCache interface.
func (r *RealRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return r.Client.Set(ctx, key, value, expiration)
}

// Get implements the Get method of the RedisCache interface.
func (r *RealRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	return r.Client.Get(ctx, key)
}

// Del implements the Del method of the RedisCache interface.
func (r *RealRedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	return r.Client.Del(ctx, keys...)
}
