package infra

import (
	"backend/internal/seo/domain"
	"context"
	"fmt"
	"log/slog"
	neturl "net/url"
	"sync"
)

type ParallelRunner struct {
	base     domain.Fetcher
	auditors []domain.Auditor
}

func NewParallelRunner(base domain.Fetcher, auditors ...domain.Auditor) *ParallelRunner {
	return &ParallelRunner{
		base:     base,
		auditors: auditors,
	}
}

func (m *ParallelRunner) Run(ctx context.Context, url *neturl.URL) (*domain.PageReport, error) {
	raw, err := m.base.Scan(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("infra: base scan failed: %w", err)
	}

	results := make([]*domain.AuditResult, len(m.auditors))
	var wg sync.WaitGroup

	for i, auditor := range m.auditors {
		wg.Add(1)
		go func(idx int, sc domain.Auditor) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					slog.Error("infra: panic in auditor",
						"auditor", sc.AuditorName(),
						"recover", r,
					)
				}
			}()

			result, err := sc.Analyze(ctx, raw)
			if err != nil {
				slog.Warn("infra: auditor returned error, skipping",
					"auditor", sc.AuditorName(),
					"error", err,
				)
				return
			}

			results[idx] = result
		}(i, auditor)
	}

	wg.Wait()

	filtered := make([]domain.AuditResult, 0, len(m.auditors))
	for _, ptr := range results {
		if ptr != nil {
			filtered = append(filtered, *ptr)
		}
	}

	return &domain.PageReport{
		URL:     url,
		Results: filtered,
	}, nil
}
