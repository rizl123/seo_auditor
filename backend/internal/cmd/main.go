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
	scanner := infrastructure.NewHttpScanner()

	redisAddr := os.Getenv("REDIS_ADDR")
	cache := infrastructure.NewRedisScannerCache(redisAddr, 24*time.Hour)

	usecase := usecase.NewScanUsecase(scanner, cache)
	handler := http.NewScanHandler(usecase)

	r := gin.Default()
	r.RedirectTrailingSlash = true

	http.SetupRouter(r, handler)

	r.Run(":8080")
}
