package shared

import (
	"context"
	"errors"
	"time"
)

var ErrCacheMiss = errors.New("cache miss")

type Cacher interface {
	Fetch(ctx context.Context, group string, key string, obj any) error
	Store(ctx context.Context, group string, key string, obj any, ttl time.Duration) error
	PingWithTimeout(timeout time.Duration) error
	Close() error
}
