package shared

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(addr string) *RedisClient {
	return &RedisClient{
		Client: redis.NewClient(&redis.Options{Addr: addr}),
	}
}

func (r *RedisClient) Fetch(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

func (r *RedisClient) Store(ctx context.Context, key string, b []byte, ttl time.Duration) error {
	return r.Client.Set(ctx, key, b, ttl).Err()
}
