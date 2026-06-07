package bootstrap

import (
	"backend/internal/config"
	"backend/internal/shared"
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	handler http.Handler
	cacher  shared.Cacher
	config  *config.Config
}

func NewApp(cfg *config.Config) *App {
	cacher := SetupCacher(cfg)
	handler := SetupHuma(cfg, cacher)

	return &App{
		handler: handler,
		cacher:  cacher,
		config:  cfg,
	}
}

func (a *App) Run() {
	srv := &http.Server{
		Addr:         ":" + a.config.AppPort,
		Handler:      a.handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	close, err := (func() (func() error, error) {
		if a.cacher != nil {
			return a.cacher.Connect(ctx)
		}
		return nil, nil
	})()
	if close == nil || err != nil {
		slog.Error("bootstrap: redis ping failed, running without cache", "error", err)
	}
	defer func() {
		err := close()
		slog.Error("Error closing cacher", "error", err)
	}()

	go func() {
		slog.Info("Server starting", "port", a.config.AppPort)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("Shutting down gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("Server exiting")
}
