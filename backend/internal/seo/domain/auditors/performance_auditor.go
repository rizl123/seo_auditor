package auditors

import (
	"backend/internal/seo/domain"
	"context"
	"strings"
	"time"
)

const (
	slowResponseThreshold = 1500 * time.Millisecond
	warnResponseThreshold = 800 * time.Millisecond
)

type PerformanceAuditor struct{}

func NewPerformanceAuditor() *PerformanceAuditor { return &PerformanceAuditor{} }

func (s *PerformanceAuditor) AuditorName() string { return "performance" }

func (s *PerformanceAuditor) Analyze(_ context.Context, report *domain.PageReport) (*domain.ScanResult, error) {
	result := &domain.ScanResult{
		AuditorName:   s.AuditorName(),
		I18nNamespace: "auditors.performance",
		Details:       []domain.Detail{},
		Problems:      []domain.Problem{},
		ScannedAt:     time.Now(),
	}

	if report.Network == nil {
		return result, nil
	}

	net := report.Network
	result.Details = append(result.Details,
		domain.Detail{
			I18nLabel: "auditors.performance.labels.response_time",
			Value:     net.ResponseTime.Milliseconds(),
			Type:      domain.DetailTypeDuration,
		},
		domain.Detail{
			I18nLabel: "auditors.performance.labels.status_code",
			Value:     report.Status,
			Type:      domain.DetailTypeBadge,
		},
		domain.Detail{
			I18nLabel: "auditors.performance.labels.server_header",
			Value:     net.Server,
			Type:      domain.DetailTypeText,
		},
		domain.Detail{
			I18nLabel: "auditors.performance.labels.content_type",
			Value:     net.ContentType,
			Type:      domain.DetailTypeText,
		},
	)

	s.checkResponseTime(result, net.ResponseTime)
	s.checkStatusAndType(result, report)

	return result, nil
}

func (s *PerformanceAuditor) checkResponseTime(result *domain.ScanResult, rt time.Duration) {
	descVars := make(map[string]any)
	descVars["ms"] = rt.Milliseconds()

	switch {
	case rt > slowResponseThreshold:
		result.Problems = append(result.Problems, domain.Problem{
			I18nNamespace:   "auditors.performance.problems.slow_ttfb",
			DescriptionVars: descVars,
			Resources: []domain.Resource{
				{Title: "web.dev: Optimize TTFB", URL: "https://web.dev/articles/optimize-ttfb"},
				{
					Title: "Google: Page speed and ranking",
					URL:   "https://developers.google.com/search/blog/2010/04/using-site-speed-in-web-search-ranking",
				},
			},
		})
	case rt > warnResponseThreshold:
		result.Problems = append(result.Problems, domain.Problem{
			I18nNamespace:   "auditors.performance.problems.approaching_threshold",
			DescriptionVars: descVars,
			Resources: []domain.Resource{
				{Title: "web.dev: Optimize TTFB", URL: "https://web.dev/articles/optimize-ttfb"},
			},
		})
	}
}

func (s *PerformanceAuditor) checkStatusAndType(result *domain.ScanResult, report *domain.PageReport) {
	if report.Status != 200 {
		descVars := make(map[string]any)
		descVars["status"] = report.Status

		result.Problems = append(result.Problems, domain.Problem{
			I18nNamespace:   "auditors.performance.problems.non_200_status",
			DescriptionVars: descVars,
			Resources: []domain.Resource{
				{
					Title: "Google: HTTP status codes",
					URL:   "https://developers.google.com/search/docs/crawling-indexing/http-network-errors",
				},
				{
					Title: "MDN: HTTP response status codes",
					URL:   "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status",
				},
			},
		})
	}

	net := report.Network
	if net.ContentType != "" && !strings.Contains(net.ContentType, "text/html") {
		descVars := make(map[string]any)
		descVars["type"] = net.ContentType

		result.Problems = append(result.Problems, domain.Problem{
			I18nNamespace:   "auditors.performance.problems.unexpected_content_type",
			DescriptionVars: descVars,
			Resources: []domain.Resource{
				{Title: "MDN: Content-Type", URL: "https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Type"},
			},
		})
	}
}
