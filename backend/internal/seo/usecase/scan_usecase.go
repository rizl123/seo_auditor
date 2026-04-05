package usecase

import (
	"backend/internal/seo/domain"
	"context"
)

type ScanUsecase struct {
	scanner domain.Scanner
}

func NewScanUsecase(s domain.Scanner) *ScanUsecase {
	return &ScanUsecase{scanner: s}
}

func (u *ScanUsecase) Execute(ctx context.Context, url string) (*domain.PageReport, error) {
	return u.scanner.Scan(ctx, url)
}
