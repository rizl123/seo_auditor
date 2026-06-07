package auditors_test

import (
	"backend/internal/seo/domain"
	"context"
	"net/url"
	"time"
)

type MockFetcher struct {
	Response *domain.RawData
	Err      error
}

func (f *MockFetcher) Scan(ctx context.Context, url *url.URL) (*domain.RawData, error) {
	return f.Response, f.Err
}

func validMeta() *domain.Metadata {
	return &domain.Metadata{
		Title:       "Valid Page Title that is more than 30 chars long",
		Description: "A perfectly fine meta description that fits in the recommended range of length for SEO purposes.",
		Canonical:   "https://example.com/page",
		OgImage:     "https://example.com/img.jpg",
		H1:          []string{"This is a single valid H1"},
	}
}

func validNetwork() *domain.NetworkInfo {
	return &domain.NetworkInfo{
		ResponseTime: 200 * time.Millisecond,
		Server:       "nginx",
		ContentType:  "text/html; charset=utf-8",
	}
}
