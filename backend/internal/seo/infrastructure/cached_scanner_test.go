package infrastructure

import (
	"backend/internal/seo/domain"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCacher struct{ mock.Mock }

func (m *MockCacher) Fetch(ctx context.Context, group string, key string, obj any) error {
	args := m.Called(ctx, group, key, obj)
	if args.Get(0) == nil {
		if val, ok := args.Get(1).(*domain.PageReport); ok && obj != nil {
			*(obj.(*domain.PageReport)) = *val
		}
		return nil
	}
	return args.Error(0)
}

func (m *MockCacher) Store(ctx context.Context, group string, key string, obj any, ttl time.Duration) error {
	return m.Called(ctx, group, key, obj, ttl).Error(0)
}

func (m *MockCacher) PingWithTimeout(d time.Duration) error { return nil }
func (m *MockCacher) Close() error                          { return nil }

type MockBaseScanner struct{ mock.Mock }

func (m *MockBaseScanner) Scan(ctx context.Context, url string) (*domain.PageReport, error) {
	args := m.Called(url)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PageReport), args.Error(1)
}

func TestCachedScanner_Scan(t *testing.T) {
	ctx := context.Background()
	targetURL := "https://example.com"

	t.Run("Should return from cache with IsCached=true", func(t *testing.T) {
		mockCacher := new(MockCacher)
		mockBase := new(MockBaseScanner)
		scanner := NewCachedScanner(mockBase, mockCacher, time.Hour)

		cachedData := &domain.PageReport{URL: targetURL, Status: 200, IsCached: false}

		mockCacher.On("Fetch", ctx, "scan", targetURL, mock.Anything).Return(nil, cachedData)

		result, err := scanner.Scan(ctx, targetURL)

		assert.NoError(t, err)
		assert.True(t, result.IsCached, "IsCached flag should be true when data is retrieved from cache")
		assert.Equal(t, targetURL, result.URL)

		mockBase.AssertNotCalled(t, "Scan", mock.Anything)
	})

	t.Run("Should scan and store with IsCached=false on cache miss", func(t *testing.T) {
		mockCacher := new(MockCacher)
		mockBase := new(MockBaseScanner)
		scanner := NewCachedScanner(mockBase, mockCacher, time.Hour)

		freshReport := &domain.PageReport{URL: targetURL, Status: 200}

		mockCacher.On("Fetch", ctx, "scan", targetURL, mock.Anything).Return(errors.New("not found"), nil)
		mockBase.On("Scan", targetURL).Return(freshReport, nil)
		mockCacher.On("Store", ctx, "scan", targetURL, freshReport, time.Hour).Return(nil)

		result, err := scanner.Scan(ctx, targetURL)

		assert.NoError(t, err)
		assert.False(t, result.IsCached, "IsCached flag should be false for a fresh scan")
		assert.Equal(t, 200, result.Status)

		mockCacher.AssertExpectations(t)
		mockBase.AssertExpectations(t)
	})
}
