package bootstrap

import (
	_ "backend/docs"

	"backend/internal/shared"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

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

func SetupRouter(handler *seoDelivery.ScanHandler) *gin.Engine {
	engine := gin.Default()
	engine.RedirectTrailingSlash = true

	api := engine.Group("/api")
	{
		api.GET("/swagger/*any", func(c *gin.Context) {
			if c.Param("any") == "/" || c.Param("any") == "" {
				c.Redirect(http.StatusMovedPermanently, "/api/swagger/index.html")
				return
			}
			ginSwagger.WrapHandler(swaggerFiles.Handler)(c)
		})

		seoDelivery.SetupRouter(api, handler)
	}

	return engine
}
