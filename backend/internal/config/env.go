package config

import (
	"log/slog"
	"os"
	"strings"
	"time"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}

func getDurationEnv(key string, fallback time.Duration) time.Duration {
	val, ok := os.LookupEnv(key)
	if !ok || val == "" {
		return fallback
	}

	d, err := time.ParseDuration(val)
	if err != nil {
		slog.Warn("Invalid duration in env, using fallback",
			slog.String("key", key),
			slog.String("value", val),
			slog.Duration("fallback", fallback))
		return fallback
	}
	return d
}

func getSliceEnv(key, fallback string) []string {
	val := getEnv(key, fallback)

	parts := strings.Split(val, ",")
	result := make([]string, 0, len(parts))

	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	if len(result) == 0 {
		return []string{"*"}
	}
	return result
}
