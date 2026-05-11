package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port            string
	Env             string // dev or prod
	DatabaseURL     string
	LogLevel        slog.Level
	ShutdownTimeout time.Duration
}

func getEnv(key string, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:        getEnv("PORT", "3000"),
		Env:         getEnv("ENV", "dev"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
	}

	if _, err := strconv.Atoi(cfg.Port); err != nil {
		return nil, fmt.Errorf("invalid PORT %q: must be a number", cfg.Port)
	}

	if cfg.Env != "dev" && cfg.Env != "prod" {
		return nil, fmt.Errorf("invalid ENV %q: must be a 'dev' or 'prod'", cfg.Env)
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	logEnv := getEnv("LOG_LEVEL", "debug")
	if err := cfg.LogLevel.UnmarshalText([]byte(logEnv)); err != nil {
		return nil, fmt.Errorf("invalid LOG_LEVEL: %w", err)
	}

	d, err := time.ParseDuration(getEnv("SHUTDOWN_TIMEOUT", "10s"))
	if err != nil {
		return nil, fmt.Errorf("invalid SHUTDOWN_TIMEOUT: %w", err)
	}
	cfg.ShutdownTimeout = d

	return cfg, nil
}
