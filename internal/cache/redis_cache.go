package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type CacheService interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string, dest interface{}) error
	Delete(ctx context.Context, key string) error
	GetOrSet(ctx context.Context, key string, dest interface{}, ttl time.Duration, fn func() (interface{}, error)) error
}

type redisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) CacheService {
	return &redisCache{
		client: client,
	}
}

func (c *redisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal cache data: %w", err)
	}

	return c.client.Set(ctx, key, data, expiration).Err()
}

func (c *redisCache) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("cache miss: %w", err)
		}
		return fmt.Errorf("failed to get cache data: %w", err)
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("failed to unmarshal cache data: %w", err)
	}

	return nil
}

func (c *redisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// GetOrSet implements a common pattern for cache-aside
func (c *redisCache) GetOrSet(ctx context.Context, key string, dest interface{}, ttl time.Duration, fn func() (interface{}, error)) error {
	// Try to get from cache first
	err := c.Get(ctx, key, dest)
	if err == nil {
		return nil
	}

	// If not in cache, execute the function
	data, err := fn()
	if err != nil {
		return fmt.Errorf("failed to execute cache loader function: %w", err)
	}

	// Store in cache
	if err := c.Set(ctx, key, data, ttl); err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	// Unmarshal into destination
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	return json.Unmarshal(dataBytes, dest)
}
