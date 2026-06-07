package infra

import (
	"backend/internal/seo/domain"
	"context"
	"net/url"
	"sync"
)

type call struct {
	wg   sync.WaitGroup
	data *domain.RawData
	err  error
}

type SingleflightFetcherProxy struct {
	base domain.Fetcher

	mu    sync.Mutex
	calls map[string]*call
}

func NewSingleflightFetcherProxy(base domain.Fetcher) *SingleflightFetcherProxy {
	return &SingleflightFetcherProxy{
		base:  base,
		calls: make(map[string]*call),
	}
}

func (s *SingleflightFetcherProxy) Scan(
	ctx context.Context,
	u *url.URL,
) (*domain.RawData, error) {
	key := u.String()

	s.mu.Lock()

	if c, ok := s.calls[key]; ok {
		s.mu.Unlock()

		c.wg.Wait()
		return c.data, c.err
	}

	c := &call{}
	c.wg.Add(1)
	s.calls[key] = c

	s.mu.Unlock()

	c.data, c.err = s.base.Scan(ctx, u)

	c.wg.Done()

	s.mu.Lock()
	delete(s.calls, key)
	s.mu.Unlock()

	return c.data, c.err
}
