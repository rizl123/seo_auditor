package shared

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
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
	fullKey := group + ":" + key
	cached, err := r.Client.Get(ctx, fullKey).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrCacheMiss
		}
		slog.Error("redis: get failed", "key", fullKey, "error", err)
		return fmt.Errorf("redis get: %w", err)
	}

	if err := json.Unmarshal([]byte(cached), &obj); err != nil {
		slog.Error("redis: unmarshal failed", "key", fullKey, "error", err)
		return fmt.Errorf("unmarshal: %w", err)
	}

	return nil
}

func (r *RedisCacher) Store(ctx context.Context, group string, key string, obj any, ttl time.Duration) error {
	fullKey := group + ":" + key
	b, err := json.Marshal(obj)

	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	err = r.Client.Set(ctx, fullKey, b, ttl).Err()
	if err != nil {
		slog.Error("redis: set failed", "key", fullKey, "error", err)
	}
	return err
}

func (r *RedisCacher) PingWithTimeout(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return r.Client.Ping(ctx).Err()
}

func (r *RedisCacher) Close() error {
	return r.Client.Close()
}
