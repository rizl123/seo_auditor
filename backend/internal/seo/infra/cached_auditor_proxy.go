package infra

import (
	"backend/internal/seo/domain"
	"backend/internal/shared"
	"context"
	"net/url"
	"time"
)

type CachedAuditorProxy struct {
	base   domain.Auditor
	cacher shared.Cacher
	ttl    time.Duration
}

func NewCachedAuditor(base domain.Auditor, cacher shared.Cacher, ttl time.Duration) *CachedAuditorProxy {
	return &CachedAuditorProxy{
		base:   base,
		cacher: cacher,
		ttl:    ttl,
	}
}

func (s *CachedAuditorProxy) AuditorName() string   { return s.base.AuditorName() }
func (s *CachedAuditorProxy) I18nNamespace() string { return s.base.I18nNamespace() }

func (s *CachedAuditorProxy) Analyze(ctx context.Context, url *url.URL) *domain.AuditResult {
	cacheKey := s.base.AuditorName() + ":" + url.String()

	var cached domain.AuditResult
	err := s.cacher.Fetch(ctx, "auditor", cacheKey, &cached)

	if err == nil {
		cached.IsCached = true
		return &cached
	}

	result := s.base.Analyze(ctx, url)
	if result.Fail != nil {
		return result
	}

	detachedCtx := context.WithoutCancel(ctx)
	_ = s.cacher.Store(detachedCtx, "auditor", cacheKey, result, s.ttl)

	result.IsCached = false
	return result
}
