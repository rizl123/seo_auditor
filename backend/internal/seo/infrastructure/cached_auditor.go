package infrastructure

import (
	"backend/internal/seo/domain"
	"backend/internal/shared"
	"context"
	"errors"
	"fmt"
	"log/slog"
	neturl "net/url"
	"sync"
	"time"
)

type CachedAuditor struct {
	base               domain.Auditor
	cacher             shared.Cacher
	ttl                time.Duration
	breakDuration      time.Duration
	cacheDisabledUntil time.Time
	mu                 sync.RWMutex
}

func NewCachedAuditor(
	base domain.Auditor,
	cacher shared.Cacher,
	ttl, breakDuration time.Duration,
) *CachedAuditor {
	return &CachedAuditor{
		base:          base,
		cacher:        cacher,
		ttl:           ttl,
		breakDuration: breakDuration,
	}
}

func (s *CachedAuditor) AuditorName() string {
	return s.base.AuditorName()
}

func (s *CachedAuditor) Analyze(ctx context.Context, report *domain.PageReport) (*domain.ScanResult, error) {
	cacheAvailable := s.isCacheAvailable()
	cacheKey := s.cacheKey(report.URL)

	if cacheAvailable {
		if result := s.fetch(ctx, cacheKey); result != nil {
			return result, nil
		}
	}

	result, err := s.base.Analyze(ctx, report)
	if err != nil {
		return nil, fmt.Errorf("infrastructure: named scanner %q failed: %w", s.AuditorName(), err)
	}

	if cacheAvailable {
		result.IsCached = false
		go s.store(ctx, cacheKey, *result)
	}

	return result, nil
}

func (s *CachedAuditor) cacheKey(u *neturl.URL) string {
	return s.base.AuditorName() + ":" + u.String()
}

func (s *CachedAuditor) fetch(ctx context.Context, key string) *domain.ScanResult {
	var result domain.ScanResult
	err := s.cacher.Fetch(ctx, "named_scan", key, &result)

	if err == nil {
		result.IsCached = true
		return &result
	}

	if !errors.Is(err, shared.ErrCacheMiss) {
		slog.Warn("infrastructure: disabling cache due to error",
			"scanner", s.AuditorName(),
			"error", err,
		)
		s.disableCache()
	}

	return nil
}

func (s *CachedAuditor) store(ctx context.Context, key string, result domain.ScanResult) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("infrastructure: panic in store goroutine",
				"recover", r,
				"scanner", s.AuditorName(),
			)
		}
	}()

	detachedCtx := context.WithoutCancel(ctx)
	bgCtx, cancel := context.WithTimeout(detachedCtx, 3*time.Second)
	defer cancel()

	err := s.cacher.Store(bgCtx, "named_scan", key, result, s.ttl)
	if err != nil {
		slog.Warn("infrastructure: failed to store in cache, tripping circuit breaker",
			"scanner", s.AuditorName(),
			"error", err,
		)
		s.disableCache()
	}
}

func (s *CachedAuditor) isCacheAvailable() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return time.Now().After(s.cacheDisabledUntil)
}

func (s *CachedAuditor) disableCache() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cacheDisabledUntil = time.Now().Add(s.breakDuration)
	// #nosec G706
	slog.Info("infrastructure: cache breaker active",
		"scanner", s.AuditorName(),
		"duration", s.breakDuration.Milliseconds(),
	)
}
