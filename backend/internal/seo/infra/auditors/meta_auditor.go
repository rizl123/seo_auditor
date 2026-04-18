package auditors

import (
	"backend/internal/seo/domain"
	"context"
	"fmt"
	"time"
)

type MetaAuditor struct{}

func NewMetaAuditor() *MetaAuditor { return &MetaAuditor{} }

func (s *MetaAuditor) AuditorName() string { return "meta" }

func (s *MetaAuditor) Analyze(_ context.Context, report *domain.PageReport) (*domain.ScanResult, error) {
	result := &domain.ScanResult{
		AuditorName: s.AuditorName(),
		Name:        "Meta & SEO Tags",
		Description: "Checks title, meta description, canonical URL, og:image " +
			"and H1 headings for correctness and completeness.",
		Details:   []domain.Detail{},
		Problems:  []domain.Problem{},
		ScannedAt: time.Now(),
	}

	if report.Metadata == nil {
		s.handleMissingMetadata(result, report.Status)
		return result, nil
	}

	meta := report.Metadata
	result.Details = append(result.Details,
		domain.Detail{Label: "Title", Value: meta.Title, Type: domain.DetailTypeText},
		domain.Detail{Label: "Description", Value: meta.Description, Type: domain.DetailTypeText},
		domain.Detail{Label: "Canonical", Value: meta.Canonical, Type: domain.DetailTypeURL},
		domain.Detail{Label: "OG Image", Value: meta.OgImage, Type: domain.DetailTypeImage},
		domain.Detail{Label: "H1 Count", Value: len(meta.H1), Type: domain.DetailTypeNumber},
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
	result.Problems = append(result.Problems, domain.Problem{
		Name: "Page metadata unavailable",
		Description: fmt.Sprintf("The page returned status %d, "+
			"so metadata could not be extracted.", status),
		Solutions: []string{"Ensure the page returns HTTP 200", "Check server logs for errors"},
		Resources: []domain.Resource{
			{Title: "HTTP Status Codes", URL: "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status"},
		},
	})
}

func (s *MetaAuditor) checkTitle(result *domain.ScanResult, title string) {
	length := len(title)
	switch {
	case title == "":
		result.Problems = append(result.Problems, domain.Problem{
			Name: "Missing title tag",
			Description: "The page has no <title> tag. Search engines use the title " +
				"as the primary clickable headline in SERPs.",
			Solutions: []string{"Add a descriptive <title> tag", "Keep it between 50–60 characters"},
			Resources: []domain.Resource{
				{Title: "Google: Title tag best practices", URL: "https://developers.google.com/search/docs/appearance/title-link"},
				{Title: "MDN: <title>", URL: "https://developer.mozilla.org/en-US/docs/Web/HTML/Element/title"},
			},
		})
	case length < 30:
		result.Problems = append(result.Problems, domain.Problem{
			Name: "Title tag too short",
			Description: fmt.Sprintf("Title is only %d characters. "+
				"Short titles look thin in SERPs.", length),
			Solutions: []string{"Expand the title to 50–60 characters"},
			Resources: []domain.Resource{
				{Title: "Moz: Title Tag", URL: "https://moz.com/learn/seo/title-tag"},
			},
		})
	case length > 60:
		result.Problems = append(result.Problems, domain.Problem{
			Name: "Title tag too long",
			Description: fmt.Sprintf("Title is %d characters — "+
				"Google typically truncates after ~60 chars.", length),
			Solutions: []string{"Trim the title to under 60 characters"},
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
			Name:        "Missing meta description",
			Description: "No meta description found. Google often uses it as the snippet in SERPs.",
			Solutions:   []string{"Add <meta name=\"description\">", "Write a summary (120–160 chars)"},
			Resources: []domain.Resource{
				{Title: "Google: Meta description", URL: "https://developers.google.com/search/docs/appearance/snippet"},
				{Title: "Ahrefs: Meta description guide", URL: "https://ahrefs.com/blog/meta-description/"},
			},
		})
	} else if length > 160 {
		result.Problems = append(result.Problems, domain.Problem{
			Name: "Meta description too long",
			Description: fmt.Sprintf("Description is %d characters. "+
				"Google truncates around 160 chars.", length),
			Solutions: []string{"Trim to 120–160 characters"},
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
			Name:        "Missing H1 heading",
			Description: "The page has no <h1> tag. H1 is a strong relevance signal.",
			Solutions:   []string{"Add exactly one <h1> that matches the topic"},
			Resources: []domain.Resource{
				{Title: "Google on headings", URL: "https://developers.google.com/search/docs/appearance/visual-elements-gallery"},
				{Title: "Moz: H1 tag", URL: "https://moz.com/learn/seo/on-page-factors"},
			},
		})
	} else if count > 1 {
		result.Problems = append(result.Problems, domain.Problem{
			Name: "Multiple H1 headings",
			Description: fmt.Sprintf("Found %d H1 tags. "+
				"This can confuse crawlers.", count),
			Solutions: []string{"Keep exactly one <h1> per page"},
			Resources: []domain.Resource{
				{Title: "Ahrefs: How many H1 tags?", URL: "https://ahrefs.com/blog/h1-tag/"},
			},
		})
	}
}

func (s *MetaAuditor) problemMissingCanonical() domain.Problem {
	return domain.Problem{
		Name:        "Missing canonical tag",
		Description: "No canonical URL specified. Search engines may index duplicate URLs.",
		Solutions:   []string{"Add <link rel=\"canonical\">"},
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
		Name:        "Missing og:image",
		Description: "No Open Graph image defined. Social shares will have no preview.",
		Solutions:   []string{"Add <meta property=\"og:image\">"},
		Resources: []domain.Resource{
			{Title: "Open Graph protocol", URL: "https://ogp.me/"},
			{Title: "Opengraph.xyz preview tool", URL: "https://www.opengraph.xyz/"},
		},
	}
}
