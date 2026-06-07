package auditors

import (
	"backend/internal/seo/domain"
	"context"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	slowResponseThreshold = 1500 * time.Millisecond
	warnResponseThreshold = 800 * time.Millisecond
)

type PerformanceAuditor struct {
	fetcher domain.Fetcher
}

func NewPerformanceAuditor(f domain.Fetcher) *PerformanceAuditor {
	return &PerformanceAuditor{fetcher: f}
}

func (s *PerformanceAuditor) AuditorName() string   { return "performance" }
func (s *PerformanceAuditor) I18nNamespace() string { return "auditors.performance" }

func (s *PerformanceAuditor) Analyze(ctx context.Context, url *url.URL) *domain.AuditResult {
	result := domain.NewAuditResult(s)

	raw, err := s.fetcher.Scan(ctx, url)

	if err != nil {
		result.Fail = &domain.AuditFail{
			Title:       "errors.auditor_failed",
			Description: "errors.try_later",
		}

		result.FinishedAt = time.Now()
		return result
	}

	if raw.Network == nil {
		result.FinishedAt = time.Now()
		return result
	}

	net := raw.Network
	result.Details = append(result.Details,
		domain.NewDurationDetail("auditors.performance.labels.response_time", net.ResponseTime),
		domain.NewBadgeDetail("auditors.performance.labels.status_code", strconv.Itoa(raw.Status)),
		domain.NewTextDetail("auditors.performance.labels.server_header", net.Server),
		domain.NewTextDetail("auditors.performance.labels.content_type", net.ContentType),
	)

	s.checkResponseTime(result, net.ResponseTime)
	s.checkStatusAndType(result, raw)

	result.FinishedAt = time.Now()
	return result
}

func (s *PerformanceAuditor) checkResponseTime(result *domain.AuditResult, rt time.Duration) {
	switch {
	case rt > slowResponseThreshold:
		problem := domain.NewProblem("auditors.performance.problems.slow_ttfb")
		problem.AddResource("web.dev: Optimize TTFB", "https://web.dev/articles/optimize-ttfb")
		problem.AddResource(
			"Google: Page speed and ranking",
			"https://developers.google.com/search/blog/2010/04/using-site-speed-in-web-search-ranking",
		)
		problem.AddInt64Var("ms", rt.Milliseconds())
		result.Problems = append(result.Problems, problem)
	case rt > warnResponseThreshold:
		problem := domain.NewProblem("auditors.performance.problems.approaching_threshold")
		problem.AddResource("web.dev: Optimize TTFB", "https://web.dev/articles/optimize-ttfb")
		problem.AddInt64Var("ms", rt.Milliseconds())
		result.Problems = append(result.Problems, problem)
	}
}

func (s *PerformanceAuditor) checkStatusAndType(result *domain.AuditResult, raw *domain.RawData) {
	if raw.Status != 200 {
		problem := domain.NewProblem("auditors.performance.problems.non_200_status")
		problem.AddResource(
			"Google: HTTP status codes",
			"https://developers.google.com/search/docs/crawling-indexing/http-network-errors",
		)
		problem.AddResource("MDN: HTTP response status codes", "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status")
		problem.AddIntVar("status", raw.Status)
		result.Problems = append(result.Problems, problem)
	}

	net := raw.Network
	if net.ContentType != "" && !strings.Contains(net.ContentType, "text/html") {
		problem := domain.NewProblem("auditors.performance.problems.unexpected_content_type")
		problem.AddResource("MDN: Content-Type", "https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Type")
		problem.AddStringVar("type", net.ContentType)
		result.Problems = append(result.Problems, problem)
	}
}
