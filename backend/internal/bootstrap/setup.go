package bootstrap

import (
	"backend/internal/shared"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"

	seoDelivery "backend/internal/seo/delivery"
	seoInfra "backend/internal/seo/infrastructure"
	seoUc "backend/internal/seo/usecase"
)

func SetupCacher() (shared.Cacher, error) {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		return nil, fmt.Errorf("REDIS_ADDR is not set")
	}

	cacher := shared.NewRedisCacher(redisAddr)

	if err := cacher.PingWithTimeout(5 * time.Second); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return cacher, nil
}

func SetupSeoHandler(cache shared.Cacher) *seoDelivery.ScanHandler {
	client := seoInfra.CreateSecureClient()
	scanner := seoInfra.NewWebScanner(client)
	reportRepo := seoInfra.NewCacheReportRepo(cache, 1*time.Hour)
	usecase := seoUc.NewScanUsecase(scanner, reportRepo)
	seoHandler := seoDelivery.NewScanHandler(usecase)
	return seoHandler
}

func SetupHuma(cacher shared.Cacher) http.Handler {
	mux := http.NewServeMux()

	config := huma.DefaultConfig("SEO Scanner API", "1.0.0")

	config.DocsPath = ""
	config.SchemasPath = "/api/schemas"
	config.OpenAPIPath = "/api/openapi"

	api := humago.New(mux, config)

	seoHandler := SetupSeoHandler(cacher)
	seoDelivery.RegisterRoutes(api, seoHandler)

	return mux
}
