package infrastructure

import (
	"backend/internal/seo/domain"
	"backend/internal/shared"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type RedisReportRepo struct {
	client *shared.RedisClient
}

func NewRedisReportRepo(client *shared.RedisClient) *RedisReportRepo {
	return &RedisReportRepo{client: client}
}

func (repo *RedisReportRepo) Fetch(ctx context.Context, url string) (*domain.PageReport, error) {
	key := fmt.Sprintf("scan:%s", url)

	cached, err := repo.client.Fetch(ctx, key)
	if err != nil {
		return nil, err
	}

	var report domain.PageReport
	if err := json.Unmarshal([]byte(cached), &report); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached report: %w", err)
	}

	return &report, nil
}

func (repo *RedisReportRepo) Store(ctx context.Context, url string, report *domain.PageReport) error {
	key := fmt.Sprintf("scan:%s", url)

	b, err := json.Marshal(report)

	if err != nil {
		return fmt.Errorf("failed to marshal report for caching: %w", err)
	}

	return repo.client.Store(ctx, key, b, 1*time.Hour)
}
