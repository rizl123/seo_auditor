package usecase

import (
	"backend/internal/seo/domain"
	"context"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRunnner struct{ mock.Mock }

func (m *MockRunnner) Run(ctx context.Context, url *url.URL) (*domain.AggregatedReport, error) {
	args := m.Called(ctx, url)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AggregatedReport), args.Error(1)
}

func TestScanUsecase_Execute(t *testing.T) {
	mockRunner := new(MockRunnner)
	uc := NewScanUsecase(mockRunner)

	targetURL, _ := url.Parse("https://test.com")

	report := &domain.AggregatedReport{
		URL: targetURL,
		Results: []domain.ScanResult{
			{AuditorName: "meta", Name: "Meta Tags"},
		},
	}

	mockRunner.On("Run", mock.Anything, targetURL).Return(report, nil)

	result, err := uc.Execute(context.Background(), targetURL)

	assert.NoError(t, err)
	assert.Equal(t, report, result)
	assert.Len(t, result.Results, 1)
	mockRunner.AssertExpectations(t)
}
