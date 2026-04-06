package config

import (
	"log/slog"
	"os"
	"time"
)

type Config struct {
	AppPort            string
	RedisAddr          string
	AllowedOrigins     []string
	CacheTTL           time.Duration
	CacheBreakDuration time.Duration
}

func Load() *Config {
	return &Config{
		AppPort:            getEnv("APP_PORT", "8080"),
		RedisAddr:          os.Getenv("REDIS_ADDR"),
		AllowedOrigins:     getSliceEnv("ALLOWED_ORIGINS", "*"),
		CacheTTL:           getDurationEnv("CACHE_TTL", 1*time.Hour),
		CacheBreakDuration: getDurationEnv("CACHE_BREAK_DURATION", 1*time.Minute),
	}
}

func (cfg *Config) Log() {
	slog.Info("Application configuration loaded",
		slog.Group("config",
			slog.String("port", cfg.AppPort),
			slog.String("redis", cfg.RedisAddr),
			slog.Any("origins", cfg.AllowedOrigins),
			slog.Duration("cache_ttl", cfg.CacheTTL),
			slog.Duration("cache_break", cfg.CacheBreakDuration),
		),
	)
}
