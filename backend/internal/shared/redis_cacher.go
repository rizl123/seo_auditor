package shared

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCacher struct {
	Client *redis.Client
}

func NewRedisCacher(addr string) *RedisCacher {
	return &RedisCacher{
		Client: redis.NewClient(&redis.Options{Addr: addr}),
	}
}

func (r *RedisCacher) Fetch(ctx context.Context, group string, key string, obj any) error {
	cached, err := r.Client.Get(ctx, group+":"+key).Result()

	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(cached), &obj); err != nil {
		return fmt.Errorf("failed to unmarshal cached report: %w", err)
	}

	return nil
}

func (r *RedisCacher) Store(ctx context.Context, group string, key string, obj any, ttl time.Duration) error {
	b, err := json.Marshal(obj)

	if err != nil {
		return fmt.Errorf("failed to marshal report for caching: %w", err)
	}

	return r.Client.Set(ctx, group+":"+key, b, ttl).Err()
}

func (r *RedisCacher) Close() error {
	return r.Client.Close()
}
