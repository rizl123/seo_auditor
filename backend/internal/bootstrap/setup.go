package bootstrap

import (
	"backend/internal/shared"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/rs/cors"

	seoDelivery "backend/internal/seo/delivery"
	"backend/internal/seo/domain"
	seoInfra "backend/internal/seo/infrastructure"
	seoUc "backend/internal/seo/usecase"
)

func SetupCacher() shared.Cacher {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		slog.Warn("bootstrap: REDIS_ADDR not set, caching disabled")
		return nil
	}

	cacher := shared.NewRedisCacher(redisAddr)
	if err := cacher.PingWithTimeout(3 * time.Second); err != nil {
		slog.Error("bootstrap: redis ping failed, running without cache", "error", err)
		return nil
	}

	return cacher
}

func SetupSeoHandler(cacher shared.Cacher) *seoDelivery.ScanHandler {
	client := seoInfra.CreateSecureClient()

	var scanner domain.Scanner = seoInfra.NewWebScanner(client)

	if cacher != nil {
		scanner = seoInfra.NewCachedScanner(scanner, cacher, 1*time.Hour, 1*time.Minute)
	} else {
		slog.Warn("bootstrap: scanner running without cache layer")
	}

	usecase := seoUc.NewScanUsecase(scanner)
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

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "*"
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   strings.Split(allowedOrigins, ","),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	return c.Handler(mux)
}
