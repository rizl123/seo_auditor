package usecase

import (
	"backend/internal/seo/domain"
	"context"
	"fmt"
	"net/url"
)

type ScanUsecase struct {
	runner domain.Runner
}

func NewScanUsecase(r domain.Runner) *ScanUsecase {
	return &ScanUsecase{runner: r}
}

func (u *ScanUsecase) Execute(ctx context.Context, url *url.URL) (*domain.AggregatedReport, error) {
	report, err := u.runner.Run(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("usecase: scan execution failed: %w", err)
	}
	return report, nil
}
