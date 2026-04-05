package usecase

import (
	"backend/internal/seo/domain"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockScanner struct{ mock.Mock }

func (m *MockScanner) Scan(ctx context.Context, url string) (*domain.PageReport, error) {
	args := m.Called(url)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PageReport), args.Error(1)
}

func TestScanUsecase_Execute(t *testing.T) {
	mockScanner := new(MockScanner)
	uc := NewScanUsecase(mockScanner)

	report := &domain.PageReport{URL: "https://test.com", Status: 200}
	mockScanner.On("Scan", "https://test.com").Return(report, nil)

	result, err := uc.Execute(context.Background(), "https://test.com")

	assert.NoError(t, err)
	assert.Equal(t, report, result)
	mockScanner.AssertExpectations(t)
}
