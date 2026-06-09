package config

import (
	"log/slog"
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppPort            int
	RedisAddr          string
	AllowedOrigins     []string
	CacheTTL           time.Duration
	CacheBreakDuration time.Duration
}

func Load() *Config {
	port, err := strconv.Atoi(getEnv("APP_PORT", "8080"))
	if err != nil {
		port = 8080
		slog.Warn("Cannot convert APP_PORT to int, using 8080")
	}
	return &Config{
		AppPort:            port,
		RedisAddr:          os.Getenv("REDIS_ADDR"),
		AllowedOrigins:     getSliceEnv("ALLOWED_ORIGINS", "*"),
		CacheTTL:           getDurationEnv("CACHE_TTL", 1*time.Hour),
		CacheBreakDuration: getDurationEnv("CACHE_BREAK_DURATION", 1*time.Minute),
	}
}

func (cfg *Config) Log() {
	slog.Info("Application configuration loaded",
		slog.Group("config",
			slog.Int("port", cfg.AppPort),
			slog.String("redis", cfg.RedisAddr),
			slog.Any("origins", cfg.AllowedOrigins),
			slog.Duration("cache_ttl", cfg.CacheTTL),
			slog.Duration("cache_break", cfg.CacheBreakDuration),
		),
	)
}
