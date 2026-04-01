package infrastructure

import (
	"backend/internal/seo/domain"
	"backend/internal/shared"
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisReportRepo_StoreAndFetch(t *testing.T) {
	mr, err := miniredis.Run()
	assert.NoError(t, err)
	defer mr.Close()

	redisClient := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	sharedClient := &shared.RedisClient{Client: redisClient}

	repo := NewRedisReportRepo(sharedClient)
	ctx := context.Background()

	report := &domain.PageReport{
		URL:    "https://example.com",
		Status: 200,
	}

	err = repo.Store(ctx, "https://example.com", report)
	assert.NoError(t, err)

	fetched, err := repo.Fetch(ctx, "https://example.com")
	assert.NoError(t, err)
	assert.NotNil(t, fetched)
	assert.Equal(t, report.URL, fetched.URL)
	assert.Equal(t, report.Status, fetched.Status)
}
