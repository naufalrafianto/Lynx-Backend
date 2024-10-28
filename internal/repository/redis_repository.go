package repository

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{client}
}

func (r *RedisRepository) StoreOTP(ctx context.Context, email, otp string) error {
	return r.client.Set(ctx, "otp:"+email, otp, 5*time.Minute).Err()
}

func (r *RedisRepository) GetOTP(ctx context.Context, email string) (string, error) {
	return r.client.Get(ctx, "otp:"+email).Result()
}
