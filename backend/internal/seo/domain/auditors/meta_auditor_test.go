package auditors_test

import (
	"context"
	"testing"

	"backend/internal/seo/domain"
	"backend/internal/seo/domain/auditors"

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
		name          string
		report        *domain.PageReport
		expectedProbs []string
		expectedVars  map[string]map[string]any
	}{
		{
			name: "Critical: Metadata is nil",
			report: &domain.PageReport{
				Status:   404,
				Metadata: nil,
			},
			expectedProbs: []string{"auditors.meta.problems.unavailable"},
			expectedVars: map[string]map[string]any{
				"auditors.meta.problems.unavailable": {"status": 404},
			},
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
			expectedProbs: []string{"auditors.meta.problems.missing_title"},
			expectedVars:  nil,
		},
		{
			name: "Title: Too short",
			report: &domain.PageReport{
				Metadata: func() *domain.Metadata {
					m := validMeta()
					m.Title = "Title that is exactly 29 chr"
					return m
				}(),
			},
			expectedProbs: []string{"auditors.meta.problems.title_too_short"},
			expectedVars: map[string]map[string]any{
				"auditors.meta.problems.title_too_short": {"length": 28},
			},
		},
		{
			name: "Title: Too long",
			report: &domain.PageReport{
				Metadata: func() *domain.Metadata {
					m := validMeta()
					m.Title = "This title is exactly sixty-one characters long for testing.."
					return m
				}(),
			},
			expectedProbs: []string{"auditors.meta.problems.title_too_long"},
			expectedVars: map[string]map[string]any{
				"auditors.meta.problems.title_too_long": {"length": 61},
			},
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
			expectedProbs: []string{"auditors.meta.problems.missing_description"},
		},
		{
			name: "Description: Too long",
			report: &domain.PageReport{
				Metadata: func() *domain.Metadata {
					m := validMeta()
					m.Description = "This is a very long description. It needs to exceed one hundred and sixty characters to trigger the auditor. So we keep writing and writing until we are absolutely sure the limit is passed."
					return m
				}(),
			},
			expectedProbs: []string{"auditors.meta.problems.description_too_long"},
			expectedVars: map[string]map[string]any{
				"auditors.meta.problems.description_too_long": {"length": 189},
			},
		},
		{
			name: "Canonical and OgImage: Missing",
			report: &domain.PageReport{
				Metadata: func() *domain.Metadata {
					m := validMeta()
					m.Canonical = ""
					m.OgImage = ""
					return m
				}(),
			},
			expectedProbs: []string{
				"auditors.meta.problems.missing_canonical",
				"auditors.meta.problems.missing_og_image",
			},
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
			expectedProbs: []string{"auditors.meta.problems.missing_h1"},
		},
		{
			name: "H1: Multiple",
			report: &domain.PageReport{
				Metadata: func() *domain.Metadata {
					m := validMeta()
					m.H1 = []string{"First", "Second", "Third"}
					return m
				}(),
			},
			expectedProbs: []string{"auditors.meta.problems.multiple_h1"},
			expectedVars: map[string]map[string]any{
				"auditors.meta.problems.multiple_h1": {"length": 3},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := auditor.Analyze(context.Background(), tt.report)
			require.NoError(t, err)

			actualProbNames := make([]string, 0, len(result.Problems))
			probMap := make(map[string]domain.Problem)

			for _, p := range result.Problems {
				actualProbNames = append(actualProbNames, p.I18nNamespace)
				probMap[p.I18nNamespace] = p
			}
			assert.ElementsMatch(t, tt.expectedProbs, actualProbNames)

			for probKey, expectedVars := range tt.expectedVars {
				actualProb, exists := probMap[probKey]
				if assert.True(t, exists) {
					for varKey, expectedVal := range expectedVars {
						assert.Equal(t, expectedVal, actualProb.DescriptionVars[varKey])
					}
				}
			}
		})
	}
}
