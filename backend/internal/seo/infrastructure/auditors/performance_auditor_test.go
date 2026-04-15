package auditors_test

import (
	"context"
	"testing"
	"time"

	"backend/internal/seo/domain"
	"backend/internal/seo/infrastructure/auditors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPerformanceAuditor_Analyze(t *testing.T) {
	auditor := auditors.NewPerformanceAuditor()

	tests := []struct {
		name          string
		report        *domain.PageReport
		expectedProbs []string
	}{
		{
			name: "Slow Response + Status Error",
			report: &domain.PageReport{
				Status: 500,
				Network: &domain.NetworkInfo{
					ResponseTime: 2000 * time.Millisecond,
					ContentType:  "text/html",
				},
			},
			expectedProbs: []string{
				"Slow server response (TTFB)",
				"Non-200 HTTP status: 500",
			},
		},
		{
			name: "Wrong Content Type Only",
			report: &domain.PageReport{
				Status: 200,
				Network: &domain.NetworkInfo{
					ResponseTime: 100 * time.Millisecond,
					ContentType:  "application/json",
				},
			},
			expectedProbs: []string{"Unexpected Content-Type"},
		},
		{
			name: "Warning threshold (801ms)",
			report: &domain.PageReport{
				Status: 200,
				Network: &domain.NetworkInfo{
					ResponseTime: 801 * time.Millisecond,
					ContentType:  "text/html",
				},
			},
			expectedProbs: []string{"Response time approaching threshold"},
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
			assert.ElementsMatch(t, tt.expectedProbs, actualProbNames)
		})
	}
}
