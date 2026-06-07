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

func (m *MockAuditor) AuditorName() string   { return "test-auditor" }
func (m *MockAuditor) I18nNamespace() string { return "auditors.test" }

func (m *MockAuditor) Analyze(ctx context.Context, targetURL *url.URL) *domain.AuditResult {
	args := m.Called(ctx, targetURL)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*domain.AuditResult)
}

func TestCachedAuditorProxy_Analyze_Logic(t *testing.T) {
	ctx := context.Background()
	targetURL, _ := url.Parse("https://example.com")
	cacheKey := "test-auditor:https://example.com"

	t.Run("CacheHit", func(t *testing.T) {
		mC, mA := new(MockCacher), new(MockAuditor)
		auditor := infra.NewCachedAuditor(mA, mC, time.Hour)

		result := domain.NewAuditResult(mA)
		result.IsCached = false // Изначально false, прокси должен выставить в true

		mC.On("Fetch", ctx, "auditor", cacheKey, mock.AnythingOfType("*domain.AuditResult")).
			Return(nil, result)

		res := auditor.Analyze(ctx, targetURL)

		assert.True(t, res.IsCached)
		mA.AssertNotCalled(t, "Analyze", mock.Anything, mock.Anything)
	})

	t.Run("CacheMiss_StoreSuccess", func(t *testing.T) {
		mC, mA := new(MockCacher), new(MockAuditor)
		auditor := infra.NewCachedAuditor(mA, mC, time.Hour)

		result := domain.NewAuditResult(mA)

		mC.On("Fetch", ctx, "auditor", cacheKey, mock.Anything).Return(shared.ErrCacheMiss)
		mA.On("Analyze", ctx, targetURL).Return(result)

		storeCalled := make(chan struct{})
		mC.On("Store", mock.Anything, "auditor", cacheKey, mock.AnythingOfType("*domain.AuditResult"), time.Hour).
			Return(nil).
			Run(func(args mock.Arguments) { close(storeCalled) })

		res := auditor.Analyze(ctx, targetURL)

		assert.False(t, res.IsCached)

		select {
		case <-storeCalled:
		case <-time.After(1 * time.Second):
			t.Fatal("timeout waiting for Store call")
		}
		mC.AssertExpectations(t)
		mA.AssertExpectations(t)
	})
}
