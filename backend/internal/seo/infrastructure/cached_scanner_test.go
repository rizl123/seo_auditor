package infrastructure

import (
	"backend/internal/seo/domain"
	"backend/internal/shared"
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
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

func (m *MockCacher) PingWithTimeout(d time.Duration) error { return m.Called(d).Error(0) }
func (m *MockCacher) Close() error                          { return m.Called().Error(0) }

type MockBaseScanner struct{ mock.Mock }

func (m *MockBaseScanner) Scan(ctx context.Context, url string) (*domain.PageReport, error) {
	args := m.Called(ctx, url)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PageReport), args.Error(1)
}

func TestCachedScanner_Scan_Logic(t *testing.T) {
	ctx := context.Background()
	targetURL := "https://example.com"
	report := &domain.PageReport{URL: targetURL, Status: 200}

	t.Run("CacheHit", func(t *testing.T) {
		mC, mB := new(MockCacher), new(MockBaseScanner)
		scanner := NewCachedScanner(mB, mC, time.Hour, time.Minute)

		mC.On("Fetch", ctx, "scan", targetURL, mock.Anything).Return(nil, report)

		res, err := scanner.Scan(ctx, targetURL)

		assert.NoError(t, err)
		assert.True(t, res.IsCached)
		mB.AssertNotCalled(t, "Scan", mock.Anything, mock.Anything)
	})

	t.Run("CacheMiss_StoreSuccess", func(t *testing.T) {
		mC, mB := new(MockCacher), new(MockBaseScanner)
		scanner := NewCachedScanner(mB, mC, time.Hour, time.Minute)

		mC.On("Fetch", ctx, "scan", targetURL, mock.Anything).Return(shared.ErrCacheMiss)
		mB.On("Scan", ctx, targetURL).Return(report, nil)

		storeCalled := make(chan struct{})
		mC.On("Store",
			mock.Anything,
			"scan",
			targetURL,
			mock.AnythingOfType("*domain.PageReport"),
			time.Hour,
		).Return(nil).Run(func(args mock.Arguments) {
			close(storeCalled)
		})

		res, err := scanner.Scan(ctx, targetURL)

		assert.NoError(t, err)
		assert.False(t, res.IsCached)

		select {
		case <-storeCalled:
		case <-time.After(1 * time.Second):
			t.Fatal("timeout waiting for Store call")
		}

		mC.AssertExpectations(t)
	})

	t.Run("CircuitBreaker_Activation_On_Fetch_Error", func(t *testing.T) {
		mC, mB := new(MockCacher), new(MockBaseScanner)
		scanner := NewCachedScanner(mB, mC, time.Hour, time.Minute)

		mC.On("Fetch", ctx, "scan", targetURL, mock.Anything).Return(assert.AnError).Once()
		mB.On("Scan", ctx, targetURL).Return(report, nil).Twice()
		mC.On("Store", mock.Anything, "scan", targetURL, mock.Anything, mock.Anything).Return(nil).Maybe()

		_, _ = scanner.Scan(ctx, targetURL)
		res, err := scanner.Scan(ctx, targetURL)

		assert.NoError(t, err)
		assert.False(t, res.IsCached)
		mC.AssertNumberOfCalls(t, "Fetch", 1)
	})

	t.Run("CircuitBreaker_Activation_On_Store_Error", func(t *testing.T) {
		mC, mB := new(MockCacher), new(MockBaseScanner)
		scanner := NewCachedScanner(mB, mC, time.Hour, time.Minute)

		mC.On("Fetch", ctx, "scan", targetURL, mock.Anything).Return(shared.ErrCacheMiss).Twice()
		mB.On("Scan", ctx, targetURL).Return(report, nil).Twice()

		storeDone := make(chan struct{})
		mC.On("Store", mock.Anything, "scan", targetURL, mock.Anything, mock.Anything).
			Return(assert.AnError).Once().Run(func(args mock.Arguments) {
			close(storeDone)
		})

		_, _ = scanner.Scan(ctx, targetURL)

		select {
		case <-storeDone:
		case <-time.After(1 * time.Second):
			t.Fatal("Store was not called")
		}

		time.Sleep(10 * time.Millisecond)

		res, err := scanner.Scan(ctx, targetURL)

		assert.NoError(t, err)
		assert.False(t, res.IsCached)
		mC.AssertNumberOfCalls(t, "Fetch", 1)
	})

	t.Run("Store_Panic_Recovery", func(t *testing.T) {
		mC, mB := new(MockCacher), new(MockBaseScanner)
		scanner := NewCachedScanner(mB, mC, time.Hour, time.Minute)

		mC.On("Fetch", ctx, "scan", targetURL, mock.Anything).Return(shared.ErrCacheMiss)
		mB.On("Scan", ctx, targetURL).Return(report, nil)

		panicDone := make(chan struct{})
		mC.On("Store", mock.Anything, "scan", targetURL, mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				defer close(panicDone)
				panic("cacher panic")
			}).Return(nil)

		_, _ = scanner.Scan(ctx, targetURL)

		select {
		case <-panicDone:
		case <-time.After(1 * time.Second):
			t.Fatal("Panic was not recovered")
		}
	})
}

func TestCachedScanner_Integration_Miniredis(t *testing.T) {
	s := miniredis.RunT(t)
	r := redis.NewClient(&redis.Options{Addr: s.Addr()})
	c := &shared.RedisCacher{Client: r}

	ctx := context.Background()
	url := "https://example.com"
	ttl := time.Minute

	mB := new(MockBaseScanner)
	scanner := NewCachedScanner(mB, c, ttl, time.Minute)
	report := &domain.PageReport{URL: url, Status: 200}

	mB.On("Scan", ctx, url).Return(report, nil).Once()

	res1, _ := scanner.Scan(ctx, url)
	assert.False(t, res1.IsCached)
	assert.Eventually(t, func() bool { return s.Exists("scan:" + url) }, 500*time.Millisecond, 10*time.Millisecond)

	res2, _ := scanner.Scan(ctx, url)
	assert.True(t, res2.IsCached)

	s.FastForward(ttl + time.Second)
	mB.On("Scan", ctx, url).Return(report, nil).Once()

	res3, _ := scanner.Scan(ctx, url)
	assert.False(t, res3.IsCached)

	mB.AssertExpectations(t)
}
