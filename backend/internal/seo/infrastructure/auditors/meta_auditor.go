package scanners

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
		Description: "Checks title, meta description, canonical URL, og:image and H1 headings for correctness and completeness.",
		Details:     map[string]any{},
		Problems:    []domain.Problem{},
		ScannedAt:   time.Now(),
	}

	if report.Metadata == nil {
		result.Problems = append(result.Problems, domain.Problem{
			Name:        "Page metadata unavailable",
			Description: fmt.Sprintf("The page returned status %d, so metadata could not be extracted.", report.Status),
			Solutions:   []string{"Ensure the page returns HTTP 200", "Check server logs for errors"},
			Resources: []domain.Resource{
				{Title: "HTTP Status Codes", URL: "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status"},
			},
		})
		return result, nil
	}

	meta := report.Metadata
	result.Details["title"] = meta.Title
	result.Details["description"] = meta.Description
	result.Details["canonical"] = meta.Canonical
	result.Details["og_image"] = meta.OgImage
	result.Details["h1_count"] = len(meta.H1)

	switch {
	case meta.Title == "":
		result.Problems = append(result.Problems, domain.Problem{
			Name:        "Missing title tag",
			Description: "The page has no <title> tag. Search engines use the title as the primary clickable headline in SERPs.",
			Solutions: []string{
				"Add a descriptive <title> tag inside <head>",
				"Keep it between 50–60 characters",
				"Include the primary keyword near the beginning",
			},
			Resources: []domain.Resource{
				{Title: "Google: Title tag best practices", URL: "https://developers.google.com/search/docs/appearance/title-link"},
				{Title: "MDN: <title>", URL: "https://developer.mozilla.org/en-US/docs/Web/HTML/Element/title"},
			},
		})
	case len(meta.Title) < 30:
		result.Problems = append(result.Problems, domain.Problem{
			Name:        "Title tag too short",
			Description: fmt.Sprintf("Title is only %d characters. Short titles miss keyword opportunities and look thin in SERPs.", len(meta.Title)),
			Solutions:   []string{"Expand the title to 50–60 characters", "Add brand name or primary keyword"},
			Resources: []domain.Resource{
				{Title: "Moz: Title Tag", URL: "https://moz.com/learn/seo/title-tag"},
			},
		})
	case len(meta.Title) > 60:
		result.Problems = append(result.Problems, domain.Problem{
			Name:        "Title tag too long",
			Description: fmt.Sprintf("Title is %d characters — Google typically truncates after ~60 chars in SERPs.", len(meta.Title)),
			Solutions:   []string{"Trim the title to under 60 characters", "Move secondary keywords to the description"},
			Resources: []domain.Resource{
				{Title: "Google: Title link documentation", URL: "https://developers.google.com/search/docs/appearance/title-link"},
			},
		})
	}

	switch {
	case meta.Description == "":
		result.Problems = append(result.Problems, domain.Problem{
			Name:        "Missing meta description",
			Description: "No meta description found. Google often uses it as the snippet in SERPs, directly affecting click-through rate.",
			Solutions: []string{
				"Add <meta name=\"description\" content=\"...\">",
				"Write a compelling summary of the page content (120–160 chars)",
				"Include a call-to-action if appropriate",
			},
			Resources: []domain.Resource{
				{Title: "Google: Meta description", URL: "https://developers.google.com/search/docs/appearance/snippet"},
				{Title: "Ahrefs: Meta description guide", URL: "https://ahrefs.com/blog/meta-description/"},
			},
		})
	case len(meta.Description) > 160:
		result.Problems = append(result.Problems, domain.Problem{
			Name:        "Meta description too long",
			Description: fmt.Sprintf("Description is %d characters. Google truncates snippets around 160 chars.", len(meta.Description)),
			Solutions:   []string{"Trim to 120–160 characters", "Put the most important info first"},
			Resources: []domain.Resource{
				{Title: "Ahrefs: Meta description length", URL: "https://ahrefs.com/blog/meta-description/"},
			},
		})
	}

	if meta.Canonical == "" {
		result.Problems = append(result.Problems, domain.Problem{
			Name:        "Missing canonical tag",
			Description: "No canonical URL specified. Without it, search engines may index duplicate URLs and split link equity.",
			Solutions: []string{
				"Add <link rel=\"canonical\" href=\"https://example.com/page\"> to <head>",
				"Use absolute URLs in canonical tags",
				"Ensure the canonical points to the preferred version of the page",
			},
			Resources: []domain.Resource{
				{Title: "Google: Canonical tag", URL: "https://developers.google.com/search/docs/crawling-indexing/consolidate-duplicate-urls"},
				{Title: "Moz: Canonical URL", URL: "https://moz.com/learn/seo/canonicalization"},
			},
		})
	}

	if meta.OgImage == "" {
		result.Problems = append(result.Problems, domain.Problem{
			Name:        "Missing og:image",
			Description: "No Open Graph image defined. Social shares of this page will have no preview image, reducing engagement.",
			Solutions: []string{
				"Add <meta property=\"og:image\" content=\"https://example.com/image.jpg\">",
				"Use an image at least 1200×630 px for best quality",
				"Also add og:title and og:description for full card support",
			},
			Resources: []domain.Resource{
				{Title: "Open Graph protocol", URL: "https://ogp.me/"},
				{Title: "Opengraph.xyz preview tool", URL: "https://www.opengraph.xyz/"},
			},
		})
	}

	switch {
	case len(meta.H1) == 0:
		result.Problems = append(result.Problems, domain.Problem{
			Name:        "Missing H1 heading",
			Description: "The page has no <h1> tag. H1 is a strong relevance signal and helps both users and crawlers understand the page topic.",
			Solutions:   []string{"Add exactly one <h1> that matches the page's main topic", "Place the H1 near the top of the content"},
			Resources: []domain.Resource{
				{Title: "Google on headings", URL: "https://developers.google.com/search/docs/appearance/visual-elements-gallery"},
				{Title: "Moz: H1 tag", URL: "https://moz.com/learn/seo/on-page-factors"},
			},
		})
	case len(meta.H1) > 1:
		result.Problems = append(result.Problems, domain.Problem{
			Name:        "Multiple H1 headings",
			Description: fmt.Sprintf("Found %d H1 tags. Multiple H1s dilute the topical signal and can confuse crawlers.", len(meta.H1)),
			Solutions:   []string{"Keep exactly one <h1> per page", "Demote additional headings to <h2> or lower"},
			Resources: []domain.Resource{
				{Title: "Ahrefs: How many H1 tags?", URL: "https://ahrefs.com/blog/h1-tag/"},
			},
		})
	}

	return result, nil
}
