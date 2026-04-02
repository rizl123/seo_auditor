package main

import (
	_ "backend/docs"
	"net/http"

	seoDelivery "backend/internal/seo/delivery"
	seoInfra "backend/internal/seo/infrastructure"
	seoUc "backend/internal/seo/usecase"
	"backend/internal/shared"

	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")
	redis := shared.NewRedisClient(redisAddr)

	r := gin.Default()
	r.RedirectTrailingSlash = true

	seoHandle := seoDelivery.NewScanHandler(
		seoUc.NewScanUsecase(
			seoInfra.NewWebScanner(
				seoInfra.CreateSecureClient(),
			),
			seoInfra.NewRedisReportRepo(redis),
		),
	)

	SetupRouter(r, seoHandle)

	r.Run(":8080")
}

func SetupRouter(r *gin.Engine, handler *seoDelivery.ScanHandler) {
	api := r.Group("/api")
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
}
