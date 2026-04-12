package infrastructure

import (
	"backend/internal/seo/domain"
	"backend/internal/shared"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

type CachedScanner struct {
	base               domain.Scanner
	cacher             shared.Cacher
	ttl                time.Duration
	breakDuration      time.Duration
	cacheDisabledUntil time.Time
	mu                 sync.RWMutex
}

func NewCachedScanner(base domain.Scanner, cacher shared.Cacher, ttl, breakDuration time.Duration) *CachedScanner {
	return &CachedScanner{
		base:          base,
		cacher:        cacher,
		ttl:           ttl,
		breakDuration: breakDuration,
	}
}

func (s *CachedScanner) Scan(ctx context.Context, url string) (*domain.PageReport, error) {
	cacheAvailable := s.isCacheAvailable()

	if cacheAvailable {
		var report domain.PageReport
		err := s.cacher.Fetch(ctx, "scan", url, &report)

		if err == nil {
			report.IsCached = true
			return &report, nil
		}

		if !errors.Is(err, shared.ErrCacheMiss) {
			slog.Warn("infrastructure: disabling cache due to error", "url", url, "error", err)
			s.disableCache()
		}
	}

	res, err := s.base.Scan(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("infrastructure: base scan failed: %w", err)
	}

	if res.Status == 200 && cacheAvailable {
		res.IsCached = false
		go s.store(ctx, url, *res)
	}

	return res, nil
}

func (s *CachedScanner) store(ctx context.Context, url string, report domain.PageReport) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("infrastructure: panic in store goroutine", "recover", r, "url", url)
		}
	}()

	detachedCtx := context.WithoutCancel(ctx)

	bgCtx, cancel := context.WithTimeout(detachedCtx, 3*time.Second)
	defer cancel()

	err := s.cacher.Store(bgCtx, "scan", url, &report, s.ttl)
	if err != nil {
		slog.Warn("infrastructure: failed to store in cache, tripping circuit breaker", "url", url, "error", err)
		s.disableCache()
	}
}

func (s *CachedScanner) isCacheAvailable() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return time.Now().After(s.cacheDisabledUntil)
}

func (s *CachedScanner) disableCache() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cacheDisabledUntil = time.Now().Add(s.breakDuration)
	slog.Info("infrastructure: cache breaker active", "duration", s.breakDuration)
}
