package delivery

import (
	"backend/internal/seo/usecase"
	"context"
	"fmt"
	"net/url"

	"golang.org/x/net/idna"

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

func normalizeURL(raw string) (*url.URL, error) {
	parsed, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}

	if parsed.Scheme == "" || parsed.Host == "" {
		return nil, fmt.Errorf("missing scheme or host")
	}

	hostASCII, err := idna.ToASCII(parsed.Host)
	if err != nil {
		return nil, err
	}

	parsed.Host = hostASCII
	return parsed, nil
}

func (h *ScanHandler) HandleScan(ctx context.Context, input *ScanInput) (*ScanOutput, error) {

	parsedURL, err := normalizeURL(input.URL)
	if err != nil {
		return nil, huma.Error400BadRequest("errors.invalid_url", fmt.Errorf("errors.url_format"))
	}

	report, err := h.usecase.Execute(ctx, parsedURL)
	if err != nil {
		return nil, huma.Error500InternalServerError("errors.internal_error", fmt.Errorf("errors.try_later"))
	}

	return &ScanOutput{Body: ToPageReportDTO(report)}, nil
}
