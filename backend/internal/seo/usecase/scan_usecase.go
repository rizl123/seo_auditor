package usecase

import (
	"backend/internal/seo/domain"
	"context"
	"net/url"
	"sync"
)

type ScanUsecase struct {
	auditors []domain.Auditor
}

func NewScanUsecase(a ...domain.Auditor) *ScanUsecase {
	return &ScanUsecase{auditors: a}
}

func (u *ScanUsecase) Execute(ctx context.Context, url *url.URL) (*domain.PageReport, error) {
	var wg sync.WaitGroup
	results := make([]*domain.AuditResult, len(u.auditors))

	for i, a := range u.auditors {
		wg.Add(1)

		go func(idx int, auditor domain.Auditor) {
			defer wg.Done()

			results[idx] = auditor.Analyze(ctx, url)
		}(i, a)
	}

	wg.Wait()

	return &domain.PageReport{
		URL:     url,
		Results: results,
	}, nil
}
