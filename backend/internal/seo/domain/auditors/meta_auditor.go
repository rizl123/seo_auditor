package auditors

import (
	"backend/internal/seo/domain"
	"context"
	"net/url"
	"time"
)

type MetaAuditor struct {
	fetcher domain.Fetcher
}

func NewMetaAuditor(f domain.Fetcher) *MetaAuditor {
	return &MetaAuditor{fetcher: f}
}

func (s *MetaAuditor) AuditorName() string   { return "meta" }
func (s *MetaAuditor) I18nNamespace() string { return "auditors.meta" }

func (s *MetaAuditor) Analyze(ctx context.Context, url *url.URL) *domain.AuditResult {
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

	if raw.Metadata == nil {
		s.handleMissingMetadata(result, raw.Status)

		result.FinishedAt = time.Now()
		return result
	}

	meta := raw.Metadata
	result.Details = append(result.Details,
		domain.NewTextDetail("auditors.meta.labels.title", meta.Title),
		domain.NewTextDetail("auditors.meta.labels.description", meta.Description),
		domain.NewURLDetail("auditors.meta.labels.canonical", meta.Canonical),
		domain.NewImageDetail("auditors.meta.labels.og_image", meta.OgImage),
		domain.NewNumberDetail("auditors.meta.labels.h1_count", len(meta.H1)),
	)

	s.checkTitle(result, meta.Title)
	s.checkDescription(result, meta.Description)
	s.checkHeadings(result, meta.H1)

	if meta.Canonical == "" {
		result.Problems = append(result.Problems, s.problemMissingCanonical())
	}
	if meta.OgImage == "" {
		result.Problems = append(result.Problems, s.problemMissingOgImage())
	}

	result.FinishedAt = time.Now()
	return result
}

func (s *MetaAuditor) handleMissingMetadata(result *domain.AuditResult, status int) {
	problem := domain.NewProblem("auditors.meta.problems.unavailable")

	problem.AddResource("HTTP Status Codes", "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status")
	problem.AddIntVar("status", status)

	result.Problems = append(result.Problems, problem)
}

func (s *MetaAuditor) checkTitle(result *domain.AuditResult, title string) {
	length := len(title)
	switch {
	case title == "":
		problem := domain.NewProblem("auditors.meta.problems.missing_title")
		problem.AddResource(
			"Google: Title tag best practices",
			"https://developers.google.com/search/docs/appearance/title-link",
		)
		problem.AddResource("MDN: <title>", "https://developer.mozilla.org/en-US/docs/Web/HTML/Element/title")
		result.Problems = append(result.Problems, problem)
	case length < 30:
		problem := domain.NewProblem("auditors.meta.problems.title_too_short")
		problem.AddResource("Moz: Title Tag", "https://moz.com/learn/seo/title-tag")
		problem.AddIntVar("length", length)
		result.Problems = append(result.Problems, problem)
	case length > 60:
		problem := domain.NewProblem("auditors.meta.problems.title_too_long")
		problem.AddResource(
			"Google: Title link documentation",
			"https://developers.google.com/search/docs/appearance/title-link",
		)
		problem.AddIntVar("length", length)
		result.Problems = append(result.Problems, problem)
	}
}

func (s *MetaAuditor) checkDescription(result *domain.AuditResult, desc string) {
	length := len(desc)

	if desc == "" {
		problem := domain.NewProblem("auditors.meta.problems.missing_description")
		problem.AddResource("Google: Meta description", "https://developers.google.com/search/docs/appearance/snippet")
		problem.AddResource("Ahrefs: Meta description guide", "https://ahrefs.com/blog/meta-description/")
		result.Problems = append(result.Problems, problem)
	} else if length > 160 {
		problem := domain.NewProblem("auditors.meta.problems.description_too_long")
		problem.AddResource("Ahrefs: Meta description length", "https://ahrefs.com/blog/meta-description/")
		problem.AddIntVar("length", length)
		result.Problems = append(result.Problems, problem)
	}
}

func (s *MetaAuditor) checkHeadings(result *domain.AuditResult, h1s []string) {
	count := len(h1s)
	if count == 0 {
		problem := domain.NewProblem("auditors.meta.problems.missing_h1")
		problem.AddResource(
			"Google on headings",
			"https://developers.google.com/search/docs/appearance/visual-elements-gallery",
		)
		problem.AddResource("Moz: H1 tag", "https://moz.com/learn/seo/on-page-factors")
		result.Problems = append(result.Problems, problem)
	} else if count > 1 {
		problem := domain.NewProblem("auditors.meta.problems.multiple_h1")
		problem.AddResource("Ahrefs: How many H1 tags?", "https://ahrefs.com/blog/h1-tag/")
		problem.AddIntVar("length", count)
		result.Problems = append(result.Problems, problem)
	}
}

func (s *MetaAuditor) problemMissingCanonical() domain.Problem {
	problem := domain.NewProblem("auditors.meta.problems.missing_canonical")
	problem.AddResource(
		"Google: Canonical tag",
		"https://developers.google.com/search/docs/crawling-indexing/consolidate-duplicate-urls",
	)
	problem.AddResource("Moz: Canonical URL", "https://moz.com/learn/seo/canonicalization")
	return problem
}

func (s *MetaAuditor) problemMissingOgImage() domain.Problem {
	problem := domain.NewProblem("auditors.meta.problems.missing_og_image")
	problem.AddResource("Open Graph protocol", "https://ogp.me/")
	problem.AddResource("Opengraph.xyz preview tool", "https://www.opengraph.xyz/")
	return problem
}
