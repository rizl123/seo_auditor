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

type MockAuditor struct{ mock.Mock }

func (m *MockAuditor) AuditorName() string   { return "meta" }
func (m *MockAuditor) I18nNamespace() string { return "auditors.meta" }
func (m *MockAuditor) Analyze(ctx context.Context, targetURL *url.URL) *domain.AuditResult {
	args := m.Called(ctx, targetURL)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*domain.AuditResult)
}

func TestScanUsecase_Execute(t *testing.T) {
	mockAuditor := new(MockAuditor)
	uc := usecase.NewScanUsecase(mockAuditor)

	targetURL, _ := url.Parse("https://test.com")

	expectedResult := &domain.AuditResult{
		AuditorName:   "meta",
		I18nNamespace: "auditors.meta",
	}

	mockAuditor.On("Analyze", mock.Anything, targetURL).Return(expectedResult)

	report, err := uc.Execute(context.Background(), targetURL)

	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.Equal(t, targetURL, report.URL)
	assert.Len(t, report.Results, 1)
	assert.Equal(t, expectedResult, report.Results[0])

	mockAuditor.AssertExpectations(t)
}
