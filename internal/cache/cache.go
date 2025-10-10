package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zercle/gofiber-skeleton/internal/metrics"
)

// Cache wraps Redis client with convenience methods
type Cache struct {
	client *redis.Client
	ctx    context.Context
}

// Config holds cache configuration
type Config struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// New creates a new cache instance
func New(config Config) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	ctx := context.Background()

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Cache{
		client: client,
		ctx:    ctx,
	}, nil
}

// Get retrieves a value from cache
func (c *Cache) Get(key string) (string, error) {
	val, err := c.client.Get(c.ctx, key).Result()
	if err == redis.Nil {
		metrics.RecordCacheAccess(false)
		return "", nil
	}
	if err != nil {
		return "", err
	}

	metrics.RecordCacheAccess(true)
	return val, nil
}

// Set stores a value in cache with expiration
func (c *Cache) Set(key string, value interface{}, expiration time.Duration) error {
	var data string

	switch v := value.(type) {
	case string:
		data = v
	default:
		jsonData, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value: %w", err)
		}
		data = string(jsonData)
	}

	return c.client.Set(c.ctx, key, data, expiration).Err()
}

// GetJSON retrieves and unmarshals a JSON value from cache
func (c *Cache) GetJSON(key string, dest interface{}) error {
	val, err := c.Get(key)
	if err != nil {
		return err
	}
	if val == "" {
		return redis.Nil
	}

	return json.Unmarshal([]byte(val), dest)
}

// SetJSON marshals and stores a value as JSON in cache
func (c *Cache) SetJSON(key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return c.Set(key, string(jsonData), expiration)
}

// Delete removes a key from cache
func (c *Cache) Delete(key string) error {
	return c.client.Del(c.ctx, key).Err()
}

// Exists checks if a key exists in cache
func (c *Cache) Exists(key string) (bool, error) {
	count, err := c.client.Exists(c.ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Expire sets expiration on an existing key
func (c *Cache) Expire(key string, expiration time.Duration) error {
	return c.client.Expire(c.ctx, key, expiration).Err()
}

// GetOrSet retrieves a value from cache, or executes function and caches result if not found
func (c *Cache) GetOrSet(key string, expiration time.Duration, fn func() (interface{}, error)) (interface{}, error) {
	// Try to get from cache first
	val, err := c.Get(key)
	if err == nil && val != "" {
		return val, nil
	}

	// Cache miss - execute function
	result, err := fn()
	if err != nil {
		return nil, err
	}

	// Store in cache
	if err := c.Set(key, result, expiration); err != nil {
		// Log error but return result anyway
		return result, nil
	}

	return result, nil
}

// Increment increments a numeric value
func (c *Cache) Increment(key string) (int64, error) {
	return c.client.Incr(c.ctx, key).Result()
}

// Decrement decrements a numeric value
func (c *Cache) Decrement(key string) (int64, error) {
	return c.client.Decr(c.ctx, key).Result()
}

// FlushAll clears all cache (use with caution!)
func (c *Cache) FlushAll() error {
	return c.client.FlushAll(c.ctx).Err()
}

// Close closes the cache connection
func (c *Cache) Close() error {
	return c.client.Close()
}

// Keys retrieves all keys matching pattern
func (c *Cache) Keys(pattern string) ([]string, error) {
	return c.client.Keys(c.ctx, pattern).Result()
}

// TTL gets the time to live for a key
func (c *Cache) TTL(key string) (time.Duration, error) {
	return c.client.TTL(c.ctx, key).Result()
}

// MGet retrieves multiple keys at once
func (c *Cache) MGet(keys ...string) ([]interface{}, error) {
	return c.client.MGet(c.ctx, keys...).Result()
}

// MSet sets multiple key-value pairs at once
func (c *Cache) MSet(pairs map[string]interface{}) error {
	return c.client.MSet(c.ctx, pairs).Err()
}
