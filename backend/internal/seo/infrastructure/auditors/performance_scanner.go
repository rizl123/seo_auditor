package scanners

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
		Description: "Analyses server response time, HTTP status code and content type to surface basic performance and configuration issues.",
		Details:     map[string]any{},
		Problems:    []domain.Problem{},
		ScannedAt:   time.Now(),
	}

	if report.Network == nil {
		return result, nil
	}

	net := report.Network
	result.Details["response_time_ms"] = net.ResponseTime.Milliseconds()
	result.Details["status"] = report.Status
	result.Details["server"] = net.Server
	result.Details["content_type"] = net.ContentType

	switch {
	case net.ResponseTime > slowResponseThreshold:
		result.Problems = append(result.Problems, domain.Problem{
			Name: "Slow server response (TTFB)",
			Description: fmt.Sprintf(
				"Server responded in %dms. Google recommends Time To First Byte under 800ms; "+
					"above 1500ms it can negatively affect rankings.",
				net.ResponseTime.Milliseconds(),
			),
			Solutions: []string{
				"Enable server-side caching (Redis, Varnish, etc.)",
				"Use a CDN to serve content closer to users",
				"Optimise database queries or heavy server-side logic",
				"Upgrade hosting resources if consistently slow",
			},
			Resources: []domain.Resource{
				{Title: "web.dev: Optimize TTFB", URL: "https://web.dev/articles/optimize-ttfb"},
				{Title: "Google: Page speed and ranking", URL: "https://developers.google.com/search/blog/2010/04/using-site-speed-in-web-search-ranking"},
			},
		})
	case net.ResponseTime > warnResponseThreshold:
		result.Problems = append(result.Problems, domain.Problem{
			Name: "Response time approaching threshold",
			Description: fmt.Sprintf(
				"Server responded in %dms. Not critical yet, but monitor — spikes can push you over the 1500ms danger zone.",
				net.ResponseTime.Milliseconds(),
			),
			Solutions: []string{
				"Profile slow endpoints and add targeted caching",
				"Check for N+1 database queries",
			},
			Resources: []domain.Resource{
				{Title: "web.dev: Optimize TTFB", URL: "https://web.dev/articles/optimize-ttfb"},
			},
		})
	}

	if report.Status != 200 {
		result.Problems = append(result.Problems, domain.Problem{
			Name:        fmt.Sprintf("Non-200 HTTP status: %d", report.Status),
			Description: fmt.Sprintf("The page returned status %d. Search engines may not index or may demote non-200 pages.", report.Status),
			Solutions: []string{
				"Ensure the canonical URL always returns 200",
				"Set up proper 301 redirects for moved content",
				"Fix 5xx errors at the server/infrastructure level",
			},
			Resources: []domain.Resource{
				{Title: "Google: HTTP status codes", URL: "https://developers.google.com/search/docs/crawling-indexing/http-network-errors"},
				{Title: "MDN: HTTP response status codes", URL: "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status"},
			},
		})
	}

	if net.ContentType != "" && !strings.Contains(net.ContentType, "text/html") {
		result.Problems = append(result.Problems, domain.Problem{
			Name:        "Unexpected Content-Type",
			Description: fmt.Sprintf("Content-Type is %q. Search engines expect text/html for regular pages.", net.ContentType),
			Solutions: []string{
				"Verify the server returns Content-Type: text/html; charset=utf-8 for HTML pages",
				"Check for misconfigured reverse proxies or middleware",
			},
			Resources: []domain.Resource{
				{Title: "MDN: Content-Type", URL: "https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Type"},
			},
		})
	}

	return result, nil
}
