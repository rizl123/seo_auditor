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

func (s *PerformanceAuditor) Analyze(_ context.Context, raw *domain.RawData) (*domain.AuditResult, error) {
	result := &domain.AuditResult{
		AuditorName:   s.AuditorName(),
		I18nNamespace: "auditors.performance",
		Details:       []domain.Detail{},
		Problems:      []domain.Problem{},
		ScannedAt:     time.Now(),
	}

	if raw.Network == nil {
		return result, nil
	}

	net := raw.Network
	result.Details = append(result.Details,
		domain.NewDetail(
			"auditors.performance.labels.response_time",
			net.ResponseTime.Milliseconds(),
			domain.DetailTypeDuration,
		),
		domain.NewDetail("auditors.performance.labels.status_code", raw.Status, domain.DetailTypeBadge),
		domain.NewDetail("auditors.performance.labels.server_header", net.Server, domain.DetailTypeText),
		domain.NewDetail("auditors.performance.labels.content_type", net.ContentType, domain.DetailTypeText),
	)

	s.checkResponseTime(result, net.ResponseTime)
	s.checkStatusAndType(result, raw)

	return result, nil
}

func (s *PerformanceAuditor) checkResponseTime(result *domain.AuditResult, rt time.Duration) {
	descVars := make(map[string]any)
	descVars["ms"] = rt.Milliseconds()

	switch {
	case rt > slowResponseThreshold:
		result.Problems = append(result.Problems, domain.Problem{
			I18nNamespace:   "auditors.performance.problems.slow_ttfb",
			DescriptionVars: descVars,
			Resources: []domain.Resource{
				domain.NewRes("web.dev: Optimize TTFB", "https://web.dev/articles/optimize-ttfb"),
				domain.NewRes(
					"Google: Page speed and ranking",
					"https://developers.google.com/search/blog/2010/04/using-site-speed-in-web-search-ranking",
				),
			},
		})
	case rt > warnResponseThreshold:
		result.Problems = append(result.Problems, domain.Problem{
			I18nNamespace:   "auditors.performance.problems.approaching_threshold",
			DescriptionVars: descVars,
			Resources: []domain.Resource{
				domain.NewRes("web.dev: Optimize TTFB", "https://web.dev/articles/optimize-ttfb"),
			},
		})
	}
}

func (s *PerformanceAuditor) checkStatusAndType(result *domain.AuditResult, raw *domain.RawData) {
	if raw.Status != 200 {
		descVars := make(map[string]any)
		descVars["status"] = raw.Status

		result.Problems = append(result.Problems, domain.Problem{
			I18nNamespace:   "auditors.performance.problems.non_200_status",
			DescriptionVars: descVars,
			Resources: []domain.Resource{
				domain.NewRes(
					"Google: HTTP status codes",
					"https://developers.google.com/search/docs/crawling-indexing/http-network-errors",
				),
				domain.NewRes("MDN: HTTP response status codes", "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status"),
			},
		})
	}

	net := raw.Network
	if net.ContentType != "" && !strings.Contains(net.ContentType, "text/html") {
		descVars := make(map[string]any)
		descVars["type"] = net.ContentType

		result.Problems = append(result.Problems, domain.Problem{
			I18nNamespace:   "auditors.performance.problems.unexpected_content_type",
			DescriptionVars: descVars,
			Resources: []domain.Resource{
				domain.NewRes("MDN: Content-Type", "https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Type"),
			},
		})
	}
}
