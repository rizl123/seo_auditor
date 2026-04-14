package infrastructure

import (
	"backend/internal/seo/domain"
	"backend/internal/shared"
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCacher struct{ mock.Mock }

func (m *MockCacher) Fetch(ctx context.Context, group string, key string, obj any) error {
	args := m.Called(ctx, group, key, obj)
	if args.Get(0) == nil {
		if val, ok := args.Get(1).(*domain.ScanResult); ok && obj != nil {
			*(obj.(*domain.ScanResult)) = *val
		}
		return nil
	}
	return args.Error(0)
}

func (m *MockCacher) Store(ctx context.Context, group string, key string, obj any, ttl time.Duration) error {
	return m.Called(ctx, group, key, obj, ttl).Error(0)
}

func (m *MockCacher) PingWithTimeout(d time.Duration) error { return m.Called(d).Error(0) }
func (m *MockCacher) Close() error                          { return m.Called().Error(0) }

type MockAuditor struct{ mock.Mock }

func (m *MockAuditor) AuditorName() string { return "test-auditor" }
func (m *MockAuditor) Analyze(ctx context.Context, report *domain.PageReport) (*domain.ScanResult, error) {
	args := m.Called(ctx, report)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ScanResult), args.Error(1)
}

func TestCachedAuditor_Analyze_Logic(t *testing.T) {
	ctx := context.Background()
	targetURL, _ := url.Parse("https://example.com")

	report := &domain.PageReport{URL: targetURL, Status: 200}
	result := &domain.ScanResult{
		AuditorName: "test-auditor",
		IsCached:    false,
		ScannedAt:   time.Now(),
	}

	cacheKey := "test-auditor:https://example.com"

	t.Run("CacheHit", func(t *testing.T) {
		mC, mA := new(MockCacher), new(MockAuditor)
		auditor := NewCachedAuditor(mA, mC, time.Hour, time.Minute)

		mC.On("Fetch", ctx, "named_scan", cacheKey, mock.AnythingOfType("*domain.ScanResult")).
			Return(nil, result)

		res, err := auditor.Analyze(ctx, report)

		assert.NoError(t, err)
		assert.True(t, res.IsCached)
		mA.AssertNotCalled(t, "Analyze", mock.Anything, mock.Anything)
	})

	t.Run("CacheMiss_StoreSuccess", func(t *testing.T) {
		mC, mA := new(MockCacher), new(MockAuditor)
		auditor := NewCachedAuditor(mA, mC, time.Hour, time.Minute)

		mC.On("Fetch", ctx, "named_scan", cacheKey, mock.Anything).Return(shared.ErrCacheMiss)
		mA.On("Analyze", ctx, report).Return(result, nil)

		storeCalled := make(chan struct{})
		mC.On("Store", mock.Anything, "named_scan", cacheKey, mock.AnythingOfType("domain.ScanResult"), time.Hour).
			Return(nil).
			Run(func(args mock.Arguments) { close(storeCalled) })

		res, err := auditor.Analyze(ctx, report)

		assert.NoError(t, err)
		assert.False(t, res.IsCached)

		select {
		case <-storeCalled:
		case <-time.After(1 * time.Second):
			t.Fatal("timeout waiting for Store call")
		}
		mC.AssertExpectations(t)
	})
}
