package usecase

import (
	"backend/internal/domain"
	"context"
	"log"
)

type SeoCache interface {
	Get(ctx context.Context, url string) (*domain.SeoData, error)
	Set(ctx context.Context, url string, data *domain.SeoData) error
}

type SeoUsecase struct {
	repo  domain.SeoRepository
	cache SeoCache
}

func NewSeoUsecase(repo domain.SeoRepository, cache SeoCache) *SeoUsecase {
	return &SeoUsecase{repo: repo, cache: cache}
}

func (u *SeoUsecase) Analyze(url string) (*domain.SeoData, error) {
	ctx := context.Background()

	if u.cache != nil {
		cached, err := u.cache.Get(ctx, url)
		if err == nil {
			return cached, nil
		}
	}

	data, err := u.repo.FetchSeoData(url)
	if err != nil {
		return nil, err
	}

	if u.cache != nil && data.Status == 200 {
		if err := u.cache.Set(ctx, url, data); err != nil {
			log.Printf("failed to cache: %v", err)
		}
	}

	return data, nil
}
