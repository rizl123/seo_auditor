package infrastructure

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
	scanners []domain.Auditor
}

func NewParallelRunner(base domain.Fetcher, scanners ...domain.Auditor) *ParallelRunner {
	return &ParallelRunner{
		base:     base,
		scanners: scanners,
	}
}

func (m *ParallelRunner) Run(ctx context.Context, url *neturl.URL) (*domain.AggregatedReport, error) {
	pageReport, err := m.base.Scan(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("infrastructure: base scan failed: %w", err)
	}

	results := make([]domain.ScanResult, len(m.scanners))
	var wg sync.WaitGroup

	for i, scanner := range m.scanners {
		wg.Add(1)
		go func(idx int, sc domain.Auditor) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					slog.Error("infrastructure: panic in named scanner",
						"scanner", sc.AuditorName(),
						"recover", r,
					)
				}
			}()

			result, err := sc.Analyze(ctx, pageReport)
			if err != nil {
				slog.Warn("infrastructure: named scanner returned error, skipping",
					"scanner", sc.AuditorName(),
					"error", err,
				)
				results[idx] = domain.ScanResult{
					AuditorName: sc.AuditorName(),
					Name:        sc.AuditorName(),
					Description: "Scanner failed to execute",
					Problems:    []domain.Problem{},
					Details:     map[string]any{"error": err.Error()},
				}
				return
			}

			results[idx] = *result
		}(i, scanner)
	}

	wg.Wait()

	return &domain.AggregatedReport{
		URL:     url,
		Results: results,
	}, nil
}
