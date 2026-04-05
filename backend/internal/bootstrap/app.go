package bootstrap

import (
	"backend/internal/shared"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	handler http.Handler
	cacher  shared.Cacher
	port    string
}

func NewApp() (*App, error) {
	cacher := SetupCacher()

	handler := SetupHuma(cacher)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	return &App{
		handler: handler,
		cacher:  cacher,
		port:    port,
	}, nil
}

func (a *App) Run() {
	srv := &http.Server{
		Addr:         ":" + a.port,
		Handler:      a.handler,
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

	if a.cacher != nil {
		if err := a.cacher.Close(); err != nil {
			log.Printf("Error closing cacher: %v", err)
		}

	}

	log.Println("Server exiting")
}
