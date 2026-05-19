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
	Body *PageReportDTO
}

func (h *ScanHandler) HandleScan(ctx context.Context, input *ScanInput) (*ScanOutput, error) {
	url, err := url.Parse(input.URL)

	if err != nil || url.Scheme == "" || url.Host == "" {
		return nil, huma.Error400BadRequest("errors.invalid_url", fmt.Errorf("errors.url_format"))
	}

	report, err := h.usecase.Execute(ctx, url)
	if err != nil {
		return nil, huma.Error500InternalServerError("errors.internal_error", fmt.Errorf("errors.try_later"))
	}

	return &ScanOutput{Body: ToPageReportDTO(report)}, nil
}
