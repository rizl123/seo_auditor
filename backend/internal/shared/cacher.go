package shared

import (
	"context"
	"errors"
	"time"
)

var ErrCacheMiss = errors.New("cache miss")

type Cacher interface {
	Connect(ctx context.Context) (func() error, error)
	Fetch(ctx context.Context, group string, key string, obj any) error
	Store(ctx context.Context, group string, key string, obj any, ttl time.Duration) error
}
