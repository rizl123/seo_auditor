package usecase

import "backend/internal/domain"

type SeoUsecase struct {
	repo domain.SeoRepository
}

func NewSeoUsecase(repo domain.SeoRepository) *SeoUsecase {
	return &SeoUsecase{repo: repo}
}

func (u *SeoUsecase) Analyze(url string) (*domain.SeoData, error) {
	// TODO добавить проверку кеша в Redis или дополнительную валидацию URL
	return u.repo.FetchSeoData(url)
}
