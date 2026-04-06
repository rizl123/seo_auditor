package bootstrap

import (
	"backend/internal/config"
	"backend/internal/shared"
	"log/slog"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/rs/cors"

	seoDelivery "backend/internal/seo/delivery"
	"backend/internal/seo/domain"
	seoInfra "backend/internal/seo/infrastructure"
	seoUc "backend/internal/seo/usecase"
)

func SetupCacher(cfg *config.Config) shared.Cacher {
	if cfg.RedisAddr == "" {
		slog.Warn("bootstrap: REDIS_ADDR not set, caching disabled")
		return nil
	}

	cacher := shared.NewRedisCacher(cfg.RedisAddr)
	if err := cacher.PingWithTimeout(3 * time.Second); err != nil {
		slog.Error("bootstrap: redis ping failed, running without cache", "error", err)
		return nil
	}

	return cacher
}

func SetupSeoHandler(cfg *config.Config, cacher shared.Cacher) *seoDelivery.ScanHandler {
	client := seoInfra.CreateSecureClient()

	var scanner domain.Scanner = seoInfra.NewWebScanner(client)

	if cacher != nil {
		scanner = seoInfra.NewCachedScanner(scanner, cacher, cfg.CacheTTL, cfg.CacheBreakDuration)
	}

	usecase := seoUc.NewScanUsecase(scanner)
	return seoDelivery.NewScanHandler(usecase)
}

func SetupHuma(cfg *config.Config, cacher shared.Cacher) http.Handler {
	mux := http.NewServeMux()
	humaConfig := huma.DefaultConfig("SEO Scanner API", "1.0.0")

	humaConfig.DocsPath = ""
	humaConfig.SchemasPath = "/api/schemas"
	humaConfig.OpenAPIPath = "/api/openapi"

	api := humago.New(mux, humaConfig)

	seoHandler := SetupSeoHandler(cfg, cacher)
	seoDelivery.RegisterRoutes(api, seoHandler)

	c := cors.New(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	return c.Handler(mux)
}
