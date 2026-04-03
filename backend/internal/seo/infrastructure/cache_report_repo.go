package infrastructure

import (
	"backend/internal/seo/domain"
	"backend/internal/shared"
	"context"
	"time"
)

type CacheReportRepo struct {
	cacher shared.Cacher
	ttl    time.Duration
}

func NewCacheReportRepo(cacher shared.Cacher, ttl time.Duration) *CacheReportRepo {
	return &CacheReportRepo{cacher: cacher, ttl: ttl}
}

func (repo *CacheReportRepo) Fetch(ctx context.Context, url string) (*domain.PageReport, error) {
	var report domain.PageReport
	err := repo.cacher.Fetch(ctx, "scan", url, &report)
	if err != nil {
		return nil, err
	}

	return &report, nil
}

func (repo *CacheReportRepo) Store(ctx context.Context, url string, report *domain.PageReport) error {
	return repo.cacher.Store(ctx, "scan", url, report, repo.ttl)
}
