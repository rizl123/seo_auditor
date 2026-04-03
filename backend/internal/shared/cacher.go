package shared

import (
	"context"
	"time"
)

type Cacher interface {
	Fetch(ctx context.Context, group string, key string, obj any) error
	Store(ctx context.Context, group string, key string, obj any, ttl time.Duration) error
	Close() error
}
