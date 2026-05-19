package infra_test

import (
	"backend/internal/seo/domain"
	"backend/internal/seo/infra"
	"backend/internal/shared"
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCacher struct{ mock.Mock }

func (m *MockCacher) Connect(ctx context.Context) (func() error, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return func() error { return nil }, args.Error(1)
	}
	return args.Get(0).(func() error), args.Error(1)
}

func (m *MockCacher) Fetch(ctx context.Context, group string, key string, obj any) error {
	args := m.Called(ctx, group, key, obj)
	if args.Get(0) == nil {
		if val, ok := args.Get(1).(*domain.AuditResult); ok && obj != nil {
			*(obj.(*domain.AuditResult)) = *val
		}
		return nil
	}
	return args.Error(0)
}

func (m *MockCacher) Store(ctx context.Context, group string, key string, obj any, ttl time.Duration) error {
	return m.Called(ctx, group, key, obj, ttl).Error(0)
}

type MockAuditor struct{ mock.Mock }

func (m *MockAuditor) AuditorName() string { return "test-auditor" }
func (m *MockAuditor) Analyze(ctx context.Context, raw *domain.RawData) (*domain.AuditResult, error) {
	args := m.Called(ctx, raw)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AuditResult), args.Error(1)
}

func TestCachedAuditorProxy_Analyze_Logic(t *testing.T) {
	ctx := context.Background()
	targetURL, _ := url.Parse("https://example.com")

	raw := &domain.RawData{URL: targetURL, Status: 200}
	result := &domain.AuditResult{
		AuditorName: "test-auditor",
		IsCached:    false,
		ScannedAt:   time.Now(),
	}

	cacheKey := "test-auditor:https://example.com"

	t.Run("CacheHit", func(t *testing.T) {
		mC, mA := new(MockCacher), new(MockAuditor)
		auditor := infra.NewCachedAuditor(mA, mC, time.Hour, time.Minute)

		mC.On("Fetch", ctx, "auditor", cacheKey, mock.AnythingOfType("*domain.AuditResult")).
			Return(nil, result)

		res, err := auditor.Analyze(ctx, raw)

		assert.NoError(t, err)
		assert.True(t, res.IsCached)
		mA.AssertNotCalled(t, "Analyze", mock.Anything, mock.Anything)
	})

	t.Run("CacheMiss_StoreSuccess", func(t *testing.T) {
		mC, mA := new(MockCacher), new(MockAuditor)
		auditor := infra.NewCachedAuditor(mA, mC, time.Hour, time.Minute)

		mC.On("Fetch", ctx, "auditor", cacheKey, mock.Anything).Return(shared.ErrCacheMiss)
		mA.On("Analyze", ctx, raw).Return(result, nil)

		storeCalled := make(chan struct{})
		mC.On("Store", mock.Anything, "auditor", cacheKey, mock.AnythingOfType("domain.AuditResult"), time.Hour).
			Return(nil).
			Run(func(args mock.Arguments) { close(storeCalled) })

		res, err := auditor.Analyze(ctx, raw)

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
