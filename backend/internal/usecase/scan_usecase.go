package usecase

import (
	"backend/internal/domain"
	"context"
)

type Cache interface {
	Fetch(ctx context.Context, url string) (*domain.PageReport, error)
	Store(ctx context.Context, url string, report *domain.PageReport) error
}

type ScanUsecase struct {
	scanner domain.ScannerRepository
	cache   Cache
}

func NewScanUsecase(s domain.ScannerRepository, c Cache) *ScanUsecase {
	return &ScanUsecase{scanner: s, cache: c}
}

func (u *ScanUsecase) Execute(url string) (*domain.PageReport, error) {
	ctx := context.Background()

	if u.cache != nil {
		if cached, err := u.cache.Fetch(ctx, url); err == nil {
			return cached, nil
		}
	}

	report, err := u.scanner.Scan(url)
	if err != nil {
		return nil, err
	}

	if u.cache != nil && report.Status == 200 {
		_ = u.cache.Store(ctx, url, report)
	}

	return report, nil
}
