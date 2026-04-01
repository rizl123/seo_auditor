package usecase

import (
	"backend/internal/seo/domain"
	"context"
)

type ScanUsecase struct {
	scanner domain.Scanner
	repo    domain.ReportRepo
}

func NewScanUsecase(s domain.Scanner, c domain.ReportRepo) *ScanUsecase {
	return &ScanUsecase{scanner: s, repo: c}
}

func (u *ScanUsecase) Execute(ctx context.Context, url string) (*domain.PageReport, error) {
	if u.repo != nil {
		if cached, err := u.repo.Fetch(ctx, url); err == nil {
			return cached, nil
		}
	}

	report, err := u.scanner.Scan(ctx, url)
	if err != nil {
		return nil, err
	}

	if u.repo != nil && report.Status == 200 {
		_ = u.repo.Store(ctx, url, report)
	}

	return report, nil
}
