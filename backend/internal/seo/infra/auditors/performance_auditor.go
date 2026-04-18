package auditors

import (
	"backend/internal/seo/domain"
	"context"
	"fmt"
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
		AuditorName: s.AuditorName(),
		Name:        "Performance & HTTP",
		Description: "Analyses server response time, HTTP status code and content type " +
			"to surface basic performance and configuration issues.",
		Details:   []domain.Detail{},
		Problems:  []domain.Problem{},
		ScannedAt: time.Now(),
	}

	if report.Network == nil {
		return result, nil
	}

	net := report.Network
	result.Details = append(result.Details,
		domain.Detail{Label: "Response Time", Value: net.ResponseTime.Milliseconds(), Type: domain.DetailTypeDuration},
		domain.Detail{Label: "Status Code", Value: report.Status, Type: domain.DetailTypeBadge},
		domain.Detail{Label: "Server Header", Value: net.Server, Type: domain.DetailTypeText},
		domain.Detail{Label: "Content Type", Value: net.ContentType, Type: domain.DetailTypeText},
	)

	s.checkResponseTime(result, net.ResponseTime)
	s.checkStatusAndType(result, report)

	return result, nil
}

func (s *PerformanceAuditor) checkResponseTime(result *domain.ScanResult, rt time.Duration) {
	ms := rt.Milliseconds()
	switch {
	case rt > slowResponseThreshold:
		result.Problems = append(result.Problems, domain.Problem{
			Name: "Slow server response (TTFB)",
			Description: fmt.Sprintf("Server responded in %dms. "+
				"Above 1500ms it can negatively affect rankings.", ms),
			Solutions: []string{"Enable server-side caching", "Use a CDN"},
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
			Name: "Response time approaching threshold",
			Description: fmt.Sprintf("Server responded in %dms. "+
				"Monitor spikes over 1500ms.", ms),
			Solutions: []string{"Profile slow endpoints", "Check for N+1 queries"},
			Resources: []domain.Resource{
				{Title: "web.dev: Optimize TTFB", URL: "https://web.dev/articles/optimize-ttfb"},
			},
		})
	}
}

func (s *PerformanceAuditor) checkStatusAndType(result *domain.ScanResult, report *domain.PageReport) {
	if report.Status != 200 {
		result.Problems = append(result.Problems, domain.Problem{
			Name: fmt.Sprintf("Non-200 HTTP status: %d", report.Status),
			Description: fmt.Sprintf("Status %d. Search engines may not index non-200 pages.",
				report.Status),
			Solutions: []string{"Ensure 200 OK", "Set up 301 redirects"},
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
		result.Problems = append(result.Problems, domain.Problem{
			Name: "Unexpected Content-Type",
			Description: fmt.Sprintf("Content-Type is %q. HTML is expected.",
				net.ContentType),
			Solutions: []string{"Verify server returns text/html"},
			Resources: []domain.Resource{
				{Title: "MDN: Content-Type", URL: "https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Type"},
			},
		})
	}
}
