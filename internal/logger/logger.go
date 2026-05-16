// 0.2: slog handler factory. dev=tint colored, prod=JSON.
// Both cmd/refresh and cmd/server call this once at boot and SetDefault the result.
package logger

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/gultekinmakif/llama-watch/internal/config"
	"github.com/lmittmann/tint"
)

func New(cfg *config.Config) (*slog.Logger, error) {
	var handler slog.Handler

	switch cfg.Env {
	case "prod":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel})

	case "dev":
		handler = tint.NewHandler(os.Stdout, &tint.Options{
			Level:      cfg.LogLevel,
			TimeFormat: time.Kitchen,
		})

	default:
		return nil, fmt.Errorf("invalid ENV %q: must be a 'dev' or 'prod'", cfg.Env)
	}

	return slog.New(handler), nil
}
