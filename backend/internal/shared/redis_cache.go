package shared

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(addr string) *RedisClient {
	return &RedisClient{
		client: redis.NewClient(&redis.Options{Addr: addr}),
	}
}

func (r *RedisClient) Fetch(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisClient) Store(ctx context.Context, key string, b []byte, ttl time.Duration) error {
	return r.client.Set(ctx, key, b, ttl).Err()
}
