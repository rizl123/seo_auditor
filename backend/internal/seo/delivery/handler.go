package delivery

import (
	"backend/internal/seo/infrastructure"
	"backend/internal/seo/usecase"
	"context"
	"fmt"
	"net/url"
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
	Body *infrastructure.PageReportDTO
}

func (h *ScanHandler) HandleScan(ctx context.Context, input *ScanInput) (*ScanOutput, error) {
	url, err := url.Parse(input.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid url: %w", err)
	}

	report, err := h.usecase.Execute(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("delivery: handle scan: %w", err)
	}

	return &ScanOutput{Body: infrastructure.NewPageReportDTO(report)}, nil
}
