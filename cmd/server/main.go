// 2.1: HTTP server entry point. Wires config, db, logger, then serves internal/server.
// Long-running. Reads what stage 1 wrote.
package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gultekinmakif/llama-watch/internal/config"
	"github.com/gultekinmakif/llama-watch/internal/db/postgres"
	"github.com/gultekinmakif/llama-watch/internal/logger"
	"github.com/gultekinmakif/llama-watch/internal/server"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	lg, err := logger.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	slog.SetDefault(lg)

	if err := postgres.New(cfg.DatabaseURL); err != nil {
		slog.Error("database connection error", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := postgres.Close(); err != nil {
			slog.Error("postgres close failed", "error", err)
		}
	}()

	if err := postgres.Migrate(); err != nil {
		slog.Error("database migration error", "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv := server.New(cfg)
	if err := srv.Start(ctx); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}
