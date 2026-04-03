package main

import (
	_ "backend/docs"
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

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
	cache := shared.NewRedisCacher(redisAddr)

	seoHandler := getSeoHandler(cache)

	engine := SetupRouter(seoHandler)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      engine,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Println("Server started on :8080")

	<-ctx.Done()
	log.Println("Shutting down gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	if err := cache.Close(); err != nil {
		log.Printf("Error closing redis: %v", err)
	}

	log.Println("Server exiting")
}

func getSeoHandler(cacher shared.Cacher) *seoDelivery.ScanHandler {
	client := seoInfra.CreateSecureClient()
	scanner := seoInfra.NewWebScanner(client)
	reportRepo := seoInfra.NewCacheReportRepo(cacher, 1*time.Hour)
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
