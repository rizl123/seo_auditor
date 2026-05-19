package usecase_test

import (
	"backend/internal/seo/domain"
	"backend/internal/seo/usecase"
	"context"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRunnner struct{ mock.Mock }

func (m *MockRunnner) Run(ctx context.Context, url *url.URL) (*domain.PageReport, error) {
	args := m.Called(ctx, url)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PageReport), args.Error(1)
}

func TestScanUsecase_Execute(t *testing.T) {
	mockRunner := new(MockRunnner)
	uc := usecase.NewScanUsecase(mockRunner)

	targetURL, _ := url.Parse("https://test.com")

	r := &domain.PageReport{
		URL: targetURL,
		Results: []domain.AuditResult{
			{AuditorName: "meta", I18nNamespace: "Meta Tags"},
		},
	}

	mockRunner.On("Run", mock.Anything, targetURL).Return(r, nil)

	report, err := uc.Execute(context.Background(), targetURL)

	assert.NoError(t, err)
	assert.Equal(t, r, report)
	assert.Len(t, report.Results, 1)
	mockRunner.AssertExpectations(t)
}
