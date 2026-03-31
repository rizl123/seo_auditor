package infrastructure

import (
	"backend/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisScannerCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisScannerCache(addr string, ttl time.Duration) *RedisScannerCache {
	return &RedisScannerCache{
		client: redis.NewClient(&redis.Options{Addr: addr}),
		ttl:    ttl,
	}
}

func (r *RedisScannerCache) Fetch(ctx context.Context, url string) (*domain.PageReport, error) {
	key := fmt.Sprintf("scan:%s", url)
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var report domain.PageReport
	if err := json.Unmarshal([]byte(val), &report); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached report: %w", err)
	}

	return &report, nil
}

func (r *RedisScannerCache) Store(ctx context.Context, url string, report *domain.PageReport) error {
	key := fmt.Sprintf("scan:%s", url)

	b, err := json.Marshal(report)
	if err != nil {
		return fmt.Errorf("failed to marshal report for caching: %w", err)
	}

	return r.client.Set(ctx, key, b, r.ttl).Err()
}
