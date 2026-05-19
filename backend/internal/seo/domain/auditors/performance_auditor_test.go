package auditors_test

import (
	"context"
	"testing"
	"time"

	"backend/internal/seo/domain"
	"backend/internal/seo/domain/auditors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validNetwork() *domain.NetworkInfo {
	return &domain.NetworkInfo{
		ResponseTime: 200 * time.Millisecond,
		Server:       "nginx",
		ContentType:  "text/html; charset=utf-8",
	}
}

func TestPerformanceAuditor_Analyze(t *testing.T) {
	auditor := auditors.NewPerformanceAuditor()

	tests := []struct {
		name          string
		raw           *domain.RawData
		expectedProbs []string
		expectedVars  map[string]map[string]any
	}{
		{
			name: "Valid 200 OK and Fast Response",
			raw: &domain.RawData{
				Status:  200,
				Network: validNetwork(),
			},
			expectedProbs: []string{},
		},
		{
			name: "Slow Response TTFB",
			raw: &domain.RawData{
				Status: 200,
				Network: func() *domain.NetworkInfo {
					n := validNetwork()
					n.ResponseTime = 1650 * time.Millisecond
					return n
				}(),
			},
			expectedProbs: []string{"auditors.performance.problems.slow_ttfb"},
			expectedVars: map[string]map[string]any{
				"auditors.performance.problems.slow_ttfb": {"ms": int64(1650)},
			},
		},
		{
			name: "Warning Response TTFB",
			raw: &domain.RawData{
				Status: 200,
				Network: func() *domain.NetworkInfo {
					n := validNetwork()
					n.ResponseTime = 850 * time.Millisecond
					return n
				}(),
			},
			expectedProbs: []string{"auditors.performance.problems.approaching_threshold"},
			expectedVars: map[string]map[string]any{
				"auditors.performance.problems.approaching_threshold": {"ms": int64(850)},
			},
		},
		{
			name: "Non-200 HTTP Status",
			raw: &domain.RawData{
				Status:  503,
				Network: validNetwork(),
			},
			expectedProbs: []string{"auditors.performance.problems.non_200_status"},
			expectedVars: map[string]map[string]any{
				"auditors.performance.problems.non_200_status": {"status": 503},
			},
		},
		{
			name: "Unexpected Content-Type",
			raw: &domain.RawData{
				Status: 200,
				Network: func() *domain.NetworkInfo {
					n := validNetwork()
					n.ContentType = "application/xml"
					return n
				}(),
			},
			expectedProbs: []string{"auditors.performance.problems.unexpected_content_type"},
			expectedVars: map[string]map[string]any{
				"auditors.performance.problems.unexpected_content_type": {"type": "application/xml"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := auditor.Analyze(context.Background(), tt.raw)
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
