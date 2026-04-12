package usecase

import (
	"backend/internal/seo/domain"
	"context"
	"fmt"
)

type ScanUsecase struct {
	scanner domain.Scanner
}

func NewScanUsecase(s domain.Scanner) *ScanUsecase {
	return &ScanUsecase{scanner: s}
}

func (u *ScanUsecase) Execute(ctx context.Context, url string) (*domain.PageReport, error) {
	report, err := u.scanner.Scan(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("usecase: scan execution failed: %w", err)
	}
	return report, nil
}
