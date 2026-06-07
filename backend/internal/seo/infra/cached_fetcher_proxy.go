package infra

import (
	"backend/internal/seo/domain"
	"backend/internal/shared"
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"
)

type CachedFetcherProxy struct {
	base   domain.Fetcher
	cacher shared.Cacher
	ttl    time.Duration
}

func NewCachedFetcherProxy(base domain.Fetcher, cacher shared.Cacher, ttl time.Duration) *CachedFetcherProxy {
	return &CachedFetcherProxy{
		base:   base,
		cacher: cacher,
		ttl:    ttl,
	}
}

func (s *CachedFetcherProxy) Scan(ctx context.Context, u *url.URL) (*domain.RawData, error) {
	key := u.String()

	var cached domain.RawData

	err := s.cacher.Fetch(ctx, "fetcher", key, &cached)
	if err == nil {
		return &cached, nil
	}

	if !errors.Is(err, shared.ErrCacheMiss) {
		data, scanErr := s.base.Scan(ctx, u)
		if scanErr != nil {
			return nil, fmt.Errorf("scan after cache fetch failure: %w", scanErr)
		}

		return data, nil
	}

	data, err := s.base.Scan(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("scan url %q: %w", u.String(), err)
	}

	detachedCtx := context.WithoutCancel(ctx)
	_ = s.cacher.Store(detachedCtx, "fetcher", key, data, s.ttl)

	return data, nil
}
