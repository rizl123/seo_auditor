package infrastructure

import (
	"backend/internal/domain"
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisSeoCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisSeoCache(addr string, ttl time.Duration) *RedisSeoCache {
	return &RedisSeoCache{
		client: redis.NewClient(&redis.Options{Addr: addr}),
		ttl:    ttl,
	}
}

func (r *RedisSeoCache) Get(ctx context.Context, url string) (*domain.SeoData, error) {
	val, err := r.client.Get(ctx, "seo:"+url).Result()
	if err != nil {
		return nil, err
	}

	var data domain.SeoData
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *RedisSeoCache) Set(ctx context.Context, url string, data *domain.SeoData) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, "seo:"+url, b, r.ttl).Err()
}
