package auditors_test

import (
	"context"
	"testing"

	"backend/internal/seo/domain"
	"backend/internal/seo/infra/auditors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validMeta() *domain.Metadata {
	return &domain.Metadata{
		Title:       "Valid Page Title that is more than 30 chars long",
		Description: "A perfectly fine meta description that fits in the recommended range of length for SEO purposes.",
		Canonical:   "https://example.com/page",
		OgImage:     "https://example.com/img.jpg",
		H1:          []string{"This is a single valid H1"},
	}
}

func TestMetaAuditor_Analyze(t *testing.T) {
	auditor := auditors.NewMetaAuditor()

	tests := []struct {
		name           string
		report         *domain.PageReport
		expectedProbs  []string
		expectedDetail map[string]any
	}{
		{
			name: "Critical: Metadata is nil",
			report: &domain.PageReport{
				Status:   404,
				Metadata: nil,
			},
			expectedProbs: []string{"Page metadata unavailable"},
		},
		{
			name: "Title: Empty string",
			report: &domain.PageReport{
				Metadata: func() *domain.Metadata {
					m := validMeta()
					m.Title = ""
					return m
				}(),
			},
			expectedProbs: []string{"Missing title tag"},
		},
		{
			name: "Title: Too short (29 chars)",
			report: &domain.PageReport{
				Metadata: func() *domain.Metadata {
					m := validMeta()
					m.Title = "Title that is exactly 29 chr"
					return m
				}(),
			},
			expectedProbs: []string{"Title tag too short"},
		},
		{
			name: "Title: Too long (61 chars)",
			report: &domain.PageReport{
				Metadata: func() *domain.Metadata {
					m := validMeta()
					m.Title = "This title is exactly sixty-one characters long for testing.."
					return m
				}(),
			},
			expectedProbs: []string{"Title tag too long"},
		},
		{
			name: "Description: Empty",
			report: &domain.PageReport{
				Metadata: func() *domain.Metadata {
					m := validMeta()
					m.Description = ""
					return m
				}(),
			},
			expectedProbs: []string{"Missing meta description"},
		},
		{
			name: "Description: Too long (>160)",
			report: &domain.PageReport{
				Metadata: func() *domain.Metadata {
					m := validMeta()
					m.Description = "This is a very long description. It needs to exceed one hundred and sixty characters to trigger the auditor. " +
						"So we keep writing and writing until we are absolutely sure the limit is passed."
					return m
				}(),
			},
			expectedProbs: []string{"Meta description too long"},
		},
		{
			name: "Canonical: Missing",
			report: &domain.PageReport{
				Metadata: func() *domain.Metadata {
					m := validMeta()
					m.Canonical = ""
					return m
				}(),
			},
			expectedProbs: []string{"Missing canonical tag"},
		},
		{
			name: "OgImage: Missing",
			report: &domain.PageReport{
				Metadata: func() *domain.Metadata {
					m := validMeta()
					m.OgImage = ""
					return m
				}(),
			},
			expectedProbs: []string{"Missing og:image"},
		},
		{
			name: "H1: Missing",
			report: &domain.PageReport{
				Metadata: func() *domain.Metadata {
					m := validMeta()
					m.H1 = []string{}
					return m
				}(),
			},
			expectedProbs: []string{"Missing H1 heading"},
		},
		{
			name: "H1: Multiple",
			report: &domain.PageReport{
				Metadata: func() *domain.Metadata {
					m := validMeta()
					m.H1 = []string{"First", "Second"}
					return m
				}(),
			},
			expectedProbs: []string{"Multiple H1 headings"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := auditor.Analyze(context.Background(), tt.report)
			require.NoError(t, err)

			actualProbNames := make([]string, 0)
			for _, p := range result.Problems {
				actualProbNames = append(actualProbNames, p.Name)
			}
			assert.ElementsMatch(t, tt.expectedProbs, actualProbNames, tt.name)
		})
	}
}
