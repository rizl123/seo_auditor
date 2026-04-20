package delivery

import (
	"backend/internal/seo/usecase"
	"context"
	"fmt"
	"net/url"

	"github.com/danielgtaylor/huma/v2"
)

type ScanHandler struct {
	usecase *usecase.ScanUsecase
}

func NewScanHandler(u *usecase.ScanUsecase) *ScanHandler {
	return &ScanHandler{usecase: u}
}

type ScanInput struct {
	URL string `query:"url" format:"uri" doc:"URL to scan" required:"true"`
}

type ScanOutput struct {
	Body *AggregatedReportDTO
}

func (h *ScanHandler) HandleScan(ctx context.Context, input *ScanInput) (*ScanOutput, error) {
	url, err := url.Parse(input.URL)

	if err != nil || url.Scheme == "" || url.Host == "" {
		return nil, huma.Error400BadRequest(
			"Invalid URL provided",
			fmt.Errorf("url must start with http:// or https:// and contain a host"),
		)
	}

	report, err := h.usecase.Execute(ctx, url)
	if err != nil {
		return nil, huma.Error500InternalServerError(
			"Failed to process scan request",
			fmt.Errorf("please try again later"),
		)
	}

	return &ScanOutput{Body: ToAggregatedReportDTO(report)}, nil
}
