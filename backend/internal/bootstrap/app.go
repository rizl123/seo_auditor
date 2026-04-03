package bootstrap

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend/internal/shared"

	"github.com/gin-gonic/gin"
)

type App struct {
	router *gin.Engine
	cacher shared.Cacher
	port   string
}

func NewApp() (*App, error) {
	cacher, err := SetupCacher()
	if err != nil {
		return nil, fmt.Errorf("failed to setup cacher: %w", err)
	}

	seoHandler := SetupSeoHandler(cacher)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	return &App{
		router: SetupRouter(seoHandler),
		cacher: cacher,
		port:   port,
	}, nil
}

func (a *App) Run() {
	srv := &http.Server{
		Addr:         ":" + a.port,
		Handler:      a.router,
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
	log.Printf("Server started on :%s", a.port)

	<-ctx.Done()
	log.Println("Shutting down gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	if err := a.cacher.Close(); err != nil {
		log.Printf("Error closing redis: %v", err)
	}

	log.Println("Server exiting")
}
