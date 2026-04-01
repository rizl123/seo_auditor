package usecase

import (
	"backend/internal/seo/domain"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockScanner struct{ mock.Mock }

func (m *MockScanner) Scan(url string) (*domain.PageReport, error) {
	args := m.Called(url)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PageReport), args.Error(1)
}

type MockRepo struct{ mock.Mock }

func (m *MockRepo) Fetch(ctx context.Context, url string) (*domain.PageReport, error) {
	args := m.Called(ctx, url)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PageReport), args.Error(1)
}
func (m *MockRepo) Store(ctx context.Context, url string, r *domain.PageReport) error {
	return m.Called(ctx, url, r).Error(0)
}

func TestScanUsecase_Execute(t *testing.T) {
	t.Run("Should return cached data if exists", func(t *testing.T) {
		mockRepo := new(MockRepo)
		mockScanner := new(MockScanner)
		uc := NewScanUsecase(mockScanner, mockRepo)

		expectedReport := &domain.PageReport{URL: "https://test.com", Status: 200}

		mockRepo.On("Fetch", mock.Anything, "https://test.com").Return(expectedReport, nil)

		result, err := uc.Execute("https://test.com")

		assert.NoError(t, err)
		assert.Equal(t, expectedReport, result)
		mockScanner.AssertNotCalled(t, "Scan", mock.Anything)
	})

	t.Run("Should scan and store if not in cache", func(t *testing.T) {
		mockRepo := new(MockRepo)
		mockScanner := new(MockScanner)
		uc := NewScanUsecase(mockScanner, mockRepo)

		report := &domain.PageReport{URL: "https://new.com", Status: 200}

		mockRepo.On("Fetch", mock.Anything, "https://new.com").Return(nil, errors.New("not found"))
		mockScanner.On("Scan", "https://new.com").Return(report, nil)
		mockRepo.On("Store", mock.Anything, "https://new.com", report).Return(nil)

		result, err := uc.Execute("https://new.com")

		assert.NoError(t, err)
		assert.Equal(t, report, result)
		mockRepo.AssertExpectations(t)
	})
}
