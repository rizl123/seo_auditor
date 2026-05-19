package auditors

import (
	"backend/internal/seo/domain"
	"context"
	"time"
)

type MetaAuditor struct{}

func NewMetaAuditor() *MetaAuditor { return &MetaAuditor{} }

func (s *MetaAuditor) AuditorName() string { return "meta" }

func (s *MetaAuditor) Analyze(_ context.Context, raw *domain.RawData) (*domain.AuditResult, error) {
	result := &domain.AuditResult{
		AuditorName:   s.AuditorName(),
		I18nNamespace: "auditors.meta",
		Details:       []domain.Detail{},
		Problems:      []domain.Problem{},
		ScannedAt:     time.Now(),
	}

	if raw.Metadata == nil {
		s.handleMissingMetadata(result, raw.Status)
		return result, nil
	}

	meta := raw.Metadata
	result.Details = append(result.Details,
		domain.NewDetail("auditors.meta.labels.title", meta.Title, domain.DetailTypeText),
		domain.NewDetail("auditors.meta.labels.description", meta.Description, domain.DetailTypeText),
		domain.NewDetail("auditors.meta.labels.canonical", meta.Canonical, domain.DetailTypeURL),
		domain.NewDetail("auditors.meta.labels.og_image", meta.OgImage, domain.DetailTypeImage),
		domain.NewDetail("auditors.meta.labels.h1_count", len(meta.H1), domain.DetailTypeNumber),
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

func (s *MetaAuditor) handleMissingMetadata(result *domain.AuditResult, status int) {
	descVars := make(map[string]any)
	descVars["status"] = status

	result.Problems = append(result.Problems, domain.Problem{
		I18nNamespace:   "auditors.meta.problems.unavailable",
		DescriptionVars: descVars,
		Resources: []domain.Resource{
			domain.NewRes("HTTP Status Codes", "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status"),
		},
	})
}

func (s *MetaAuditor) checkTitle(result *domain.AuditResult, title string) {
	length := len(title)

	descVars := make(map[string]any)
	descVars["length"] = length

	switch {
	case title == "":
		result.Problems = append(result.Problems, domain.Problem{
			I18nNamespace: "auditors.meta.problems.missing_title",
			Resources: []domain.Resource{
				domain.NewRes(
					"Google: Title tag best practices",
					"https://developers.google.com/search/docs/appearance/title-link",
				),
				domain.NewRes("MDN: <title>", "https://developer.mozilla.org/en-US/docs/Web/HTML/Element/title"),
			},
		})
	case length < 30:
		result.Problems = append(result.Problems, domain.Problem{
			I18nNamespace:   "auditors.meta.problems.title_too_short",
			DescriptionVars: descVars,
			Resources: []domain.Resource{
				domain.NewRes("Moz: Title Tag", "https://moz.com/learn/seo/title-tag"),
			},
		})
	case length > 60:
		result.Problems = append(result.Problems, domain.Problem{
			I18nNamespace:   "auditors.meta.problems.title_too_long",
			DescriptionVars: descVars,
			Resources: []domain.Resource{
				domain.NewRes(
					"Google: Title link documentation",
					"https://developers.google.com/search/docs/appearance/title-link",
				),
			},
		})
	}
}

func (s *MetaAuditor) checkDescription(result *domain.AuditResult, desc string) {
	length := len(desc)

	if desc == "" {
		result.Problems = append(result.Problems, domain.Problem{
			I18nNamespace: "auditors.meta.problems.missing_description",
			Resources: []domain.Resource{
				domain.NewRes("Google: Meta description", "https://developers.google.com/search/docs/appearance/snippet"),
				domain.NewRes("Ahrefs: Meta description guide", "https://ahrefs.com/blog/meta-description/"),
			},
		})
	} else if length > 160 {
		descVars := make(map[string]any)
		descVars["length"] = length

		result.Problems = append(result.Problems, domain.Problem{
			I18nNamespace:   "auditors.meta.problems.description_too_long",
			DescriptionVars: descVars,
			Resources: []domain.Resource{
				domain.NewRes("Ahrefs: Meta description length", "https://ahrefs.com/blog/meta-description/"),
			},
		})
	}
}

func (s *MetaAuditor) checkHeadings(result *domain.AuditResult, h1s []string) {
	count := len(h1s)
	if count == 0 {
		result.Problems = append(result.Problems, domain.Problem{
			I18nNamespace: "auditors.meta.problems.missing_h1",
			Resources: []domain.Resource{
				domain.NewRes("Google on headings", "https://developers.google.com/search/docs/appearance/visual-elements-gallery"),
				domain.NewRes("Moz: H1 tag", "https://moz.com/learn/seo/on-page-factors"),
			},
		})
	} else if count > 1 {
		descVars := make(map[string]any)
		descVars["length"] = count

		result.Problems = append(result.Problems, domain.Problem{
			I18nNamespace:   "auditors.meta.problems.multiple_h1",
			DescriptionVars: descVars,
			Resources: []domain.Resource{
				domain.NewRes("Ahrefs: How many H1 tags?", "https://ahrefs.com/blog/h1-tag/"),
			},
		})
	}
}

func (s *MetaAuditor) problemMissingCanonical() domain.Problem {
	return domain.Problem{
		I18nNamespace: "auditors.meta.problems.missing_canonical",
		Resources: []domain.Resource{
			domain.NewRes(
				"Google: Canonical tag",
				"https://developers.google.com/search/docs/crawling-indexing/consolidate-duplicate-urls",
			),
			domain.NewRes("Moz: Canonical URL", "https://moz.com/learn/seo/canonicalization"),
		},
	}
}

func (s *MetaAuditor) problemMissingOgImage() domain.Problem {
	return domain.Problem{
		I18nNamespace: "auditors.meta.problems.missing_og_image",
		Resources: []domain.Resource{
			domain.NewRes("Open Graph protocol", "https://ogp.me/"),
			domain.NewRes("Opengraph.xyz preview tool", "https://www.opengraph.xyz/"),
		},
	}
}
