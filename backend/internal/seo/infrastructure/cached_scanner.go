package infrastructure

import (
	"backend/internal/seo/domain"
	"backend/internal/shared"
	"context"
	"time"
)

type CachedScanner struct {
	base   domain.Scanner
	cacher shared.Cacher
	ttl    time.Duration
}

func NewCachedScanner(base domain.Scanner, cacher shared.Cacher, ttl time.Duration) *CachedScanner {
	return &CachedScanner{
		base:   base,
		cacher: cacher,
		ttl:    ttl,
	}
}

func (s *CachedScanner) Scan(ctx context.Context, url string) (*domain.PageReport, error) {
	var report domain.PageReport

	err := s.cacher.Fetch(ctx, "scan", url, &report)
	if err == nil {
		report.IsCached = true
		return &report, nil
	}

	res, err := s.base.Scan(ctx, url)
	if err != nil {
		return nil, err
	}

	if res.Status == 200 {
		res.IsCached = false
		_ = s.cacher.Store(ctx, "scan", url, res, s.ttl)
	}

	return res, nil
}
