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
	addr             string
	firstPingTimeout time.Duration
	client           *redis.Client
}

func NewRedisCacher(addr string, firstPingTimeout time.Duration) *RedisCacher {
	return &RedisCacher{
		addr:             addr,
		firstPingTimeout: firstPingTimeout,
	}
}

func (r *RedisCacher) Connect(ctx context.Context) (func() error, error) {
	r.client = redis.NewClient(&redis.Options{Addr: r.addr})
	timeoutCtx, cancel := context.WithTimeout(ctx, r.firstPingTimeout)
	defer cancel()

	err := r.client.Ping(timeoutCtx).Err()
	if err != nil {
		return nil, fmt.Errorf("redis: ping failed: %w", err)
	}

	close := func() error {
		if r.client != nil {
			return r.client.Close()
		}

		return nil
	}

	return close, nil
}

func (r *RedisCacher) Fetch(ctx context.Context, group string, key string, obj any) error {
	fullKey := group + ":" + key
	cached, err := r.client.Get(ctx, fullKey).Result()

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

	err = r.client.Set(ctx, fullKey, b, ttl).Err()
	if err != nil {
		slog.Error("redis: set failed", "key", fullKey, "error", err)
		return fmt.Errorf("redis: store failed: %w", err)
	}

	return nil
}
