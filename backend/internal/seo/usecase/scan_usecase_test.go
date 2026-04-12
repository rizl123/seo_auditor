package usecase

import (
	"backend/internal/seo/domain"
	"context"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockScanner struct{ mock.Mock }

func (m *MockScanner) Scan(ctx context.Context, url *url.URL) (*domain.PageReport, error) {
	args := m.Called(url)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PageReport), args.Error(1)
}

func TestScanUsecase_Execute(t *testing.T) {
	mockScanner := new(MockScanner)
	uc := NewScanUsecase(mockScanner)

	url, _ := url.Parse("https://test.com")

	report := &domain.PageReport{URL: url, Status: 200}
	mockScanner.On("Scan", url).Return(report, nil)

	result, err := uc.Execute(context.Background(), url)

	assert.NoError(t, err)
	assert.Equal(t, report, result)
	mockScanner.AssertExpectations(t)
}
