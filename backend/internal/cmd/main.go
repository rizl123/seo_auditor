package main

import (
	"backend/internal/delivery/http"
	"backend/internal/infrastructure"
	"backend/internal/usecase"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	repo := infrastructure.NewHttpSeoRepository()

	redisAddr := os.Getenv("REDIS_ADDR")
	cache := infrastructure.NewRedisSeoCache(redisAddr, 24*time.Hour)

	seoUsecase := usecase.NewSeoUsecase(repo, cache)
	handler := http.NewSeoHandler(seoUsecase)

	r := gin.Default()
	r.RedirectTrailingSlash = true

	http.SetupRouter(r, handler)

	r.Run(":8080")
}
