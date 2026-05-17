package auditors

import (
	"backend/internal/seo/domain"
	"context"
	"time"
)

type MetaAuditor struct{}

func NewMetaAuditor() *MetaAuditor { return &MetaAuditor{} }

func (s *MetaAuditor) AuditorName() string { return "meta" }

func (s *MetaAuditor) Analyze(_ context.Context, report *domain.PageReport) (*domain.ScanResult, error) {
	result := &domain.ScanResult{
		AuditorName:   s.AuditorName(),
		I18nNamespace: "auditors.meta",
		Details:       []domain.Detail{},
		Problems:      []domain.Problem{},
		ScannedAt:     time.Now(),
	}

	if report.Metadata == nil {
		s.handleMissingMetadata(result, report.Status)
		return result, nil
	}

	meta := report.Metadata
	result.Details = append(result.Details,
		domain.Detail{I18nLabel: "auditors.meta.labels.title", Value: meta.Title, Type: domain.DetailTypeText},
		domain.Detail{I18nLabel: "auditors.meta.labels.description", Value: meta.Description, Type: domain.DetailTypeText},
		domain.Detail{I18nLabel: "auditors.meta.labels.canonical", Value: meta.Canonical, Type: domain.DetailTypeURL},
		domain.Detail{I18nLabel: "auditors.meta.labels.og_image", Value: meta.OgImage, Type: domain.DetailTypeImage},
		domain.Detail{I18nLabel: "auditors.meta.labels.h1_count", Value: len(meta.H1), Type: domain.DetailTypeNumber},
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

	return result, nil
}

func (s *MetaAuditor) handleMissingMetadata(result *domain.ScanResult, status int) {
	descVars := make(map[string]any)
	descVars["status"] = status

	result.Problems = append(result.Problems, domain.Problem{
		I18nNamespace:   "auditors.meta.problems.unavailable",
		DescriptionVars: descVars,
		Resources: []domain.Resource{
			{Title: "HTTP Status Codes", URL: "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status"},
		},
	})
}

func (s *MetaAuditor) checkTitle(result *domain.ScanResult, title string) {
	length := len(title)

	descVars := make(map[string]any)
	descVars["length"] = length

	switch {
	case title == "":
		result.Problems = append(result.Problems, domain.Problem{
			I18nNamespace: "auditors.meta.problems.missing_title",
			Resources: []domain.Resource{
				{Title: "Google: Title tag best practices", URL: "https://developers.google.com/search/docs/appearance/title-link"},
				{Title: "MDN: <title>", URL: "https://developer.mozilla.org/en-US/docs/Web/HTML/Element/title"},
			},
		})
	case length < 30:
		result.Problems = append(result.Problems, domain.Problem{
			I18nNamespace:   "auditors.meta.problems.title_too_short",
			DescriptionVars: descVars,
			Resources: []domain.Resource{
				{Title: "Moz: Title Tag", URL: "https://moz.com/learn/seo/title-tag"},
			},
		})
	case length > 60:
		result.Problems = append(result.Problems, domain.Problem{
			I18nNamespace:   "auditors.meta.problems.title_too_long",
			DescriptionVars: descVars,
			Resources: []domain.Resource{
				{Title: "Google: Title link documentation", URL: "https://developers.google.com/search/docs/appearance/title-link"},
			},
		})
	}
}

func (s *MetaAuditor) checkDescription(result *domain.ScanResult, desc string) {
	length := len(desc)

	if desc == "" {
		result.Problems = append(result.Problems, domain.Problem{
			I18nNamespace: "auditors.meta.problems.missing_description",
			Resources: []domain.Resource{
				{Title: "Google: Meta description", URL: "https://developers.google.com/search/docs/appearance/snippet"},
				{Title: "Ahrefs: Meta description guide", URL: "https://ahrefs.com/blog/meta-description/"},
			},
		})
	} else if length > 160 {
		descVars := make(map[string]any)
		descVars["length"] = length

		result.Problems = append(result.Problems, domain.Problem{
			I18nNamespace:   "auditors.meta.problems.description_too_long",
			DescriptionVars: descVars,
			Resources: []domain.Resource{
				{Title: "Ahrefs: Meta description length", URL: "https://ahrefs.com/blog/meta-description/"},
			},
		})
	}
}

func (s *MetaAuditor) checkHeadings(result *domain.ScanResult, h1s []string) {
	count := len(h1s)
	if count == 0 {
		result.Problems = append(result.Problems, domain.Problem{
			I18nNamespace: "auditors.meta.problems.missing_h1",
			Resources: []domain.Resource{
				{Title: "Google on headings", URL: "https://developers.google.com/search/docs/appearance/visual-elements-gallery"},
				{Title: "Moz: H1 tag", URL: "https://moz.com/learn/seo/on-page-factors"},
			},
		})
	} else if count > 1 {
		descVars := make(map[string]any)
		descVars["length"] = count

		result.Problems = append(result.Problems, domain.Problem{
			I18nNamespace:   "auditors.meta.problems.multiple_h1",
			DescriptionVars: descVars,
			Resources: []domain.Resource{
				{Title: "Ahrefs: How many H1 tags?", URL: "https://ahrefs.com/blog/h1-tag/"},
			},
		})
	}
}

func (s *MetaAuditor) problemMissingCanonical() domain.Problem {
	return domain.Problem{
		I18nNamespace: "auditors.meta.problems.missing_canonical",
		Resources: []domain.Resource{
			{
				Title: "Google: Canonical tag",
				URL:   "https://developers.google.com/search/docs/crawling-indexing/consolidate-duplicate-urls",
			},
			{Title: "Moz: Canonical URL", URL: "https://moz.com/learn/seo/canonicalization"},
		},
	}
}

func (s *MetaAuditor) problemMissingOgImage() domain.Problem {
	return domain.Problem{
		I18nNamespace: "auditors.meta.problems.missing_og_image",
		Resources: []domain.Resource{
			{Title: "Open Graph protocol", URL: "https://ogp.me/"},
			{Title: "Opengraph.xyz preview tool", URL: "https://www.opengraph.xyz/"},
		},
	}
}
