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

func SetupSeoHandler(cacher shared.Cacher) *seoDelivery.ScanHandler {
	client := seoInfra.CreateSecureClient()
	baseScanner := seoInfra.NewWebScanner(client)
	cachedScanner := seoInfra.NewCachedScanner(baseScanner, cacher, 1*time.Hour)
	usecase := seoUc.NewScanUsecase(cachedScanner)

	return seoDelivery.NewScanHandler(usecase)
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
