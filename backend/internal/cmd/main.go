package main

import (
	_ "backend/docs"
	"net/http"

	"backend/internal/seo/delivery"
	"backend/internal/seo/infrastructure"
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

	SetupRouter(r, redis)

	r.Run(":8080")
}

func SetupRouter(r *gin.Engine, redis *shared.RedisClient) {
	api := r.Group("/api")
	{
		api.GET("/swagger/*any", func(c *gin.Context) {
			if c.Param("any") == "/" || c.Param("any") == "" {
				c.Redirect(http.StatusMovedPermanently, "/api/swagger/index.html")
				return
			}
			ginSwagger.WrapHandler(swaggerFiles.Handler)(c)
		})

		delivery.SetupRouter(api, infrastructure.NewRedisReportRepo(redis))
	}
}
