package delivery

import (
	"backend/internal/seo/domain"
	"backend/internal/seo/usecase"
	"context"
	"fmt"
)

type ScanHandler struct {
	usecase *usecase.ScanUsecase
}

func NewScanHandler(u *usecase.ScanUsecase) *ScanHandler {
	return &ScanHandler{usecase: u}
}

type ScanInput struct {
	URL string `query:"url" json:"url" doc:"URL to scan" required:"true"`
}

type ScanOutput struct {
	Body *domain.PageReport
}

func (h *ScanHandler) HandleScan(ctx context.Context, input *ScanInput) (*ScanOutput, error) {
	report, err := h.usecase.Execute(ctx, input.URL)
	if err != nil {
		return nil, fmt.Errorf("delivery: handle scan: %w", err)
	}

	return &ScanOutput{Body: report}, nil
}
