// 1.1: One-shot orchestrator. Invoked by scripts/refresh.sh.
// Runs once per cron tick: gate by recent refresh_run, then walk + load + build, then record outcome.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	"github.com/gultekinmakif/llama-watch/internal/config"
	"github.com/gultekinmakif/llama-watch/internal/db/postgres"
	"github.com/gultekinmakif/llama-watch/internal/dimensions"
	"github.com/gultekinmakif/llama-watch/internal/logger"
	"github.com/gultekinmakif/llama-watch/internal/models"
	"gorm.io/gorm"
)

func main() {
	intervalSec := flag.Int("interval", 3300, "seconds; skip if the last refresh finished within this window")
	upstreamDir := flag.String("upstream-dir", "var/upstream", "parent of the cloned upstream repos")
	protocolsJSON := flag.String("protocols-json", "var/extracted/protocols.json", "path to the bun extractor output")
	flag.String("snapshot-out", "var/snapshot/snapshot.json", "snapshot writer output path; wired in a later step")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	lg, err := logger.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	slog.SetDefault(lg)

	lg.Info("flags parsed",
		"interval", *intervalSec,
		"upstream_dir", *upstreamDir,
		"protocols_json", *protocolsJSON,
	)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := postgres.New(cfg.DatabaseURL); err != nil {
		slog.Error("database connection error", "error", err)
		return
	}
	defer func() {
		if err := postgres.Close(); err != nil {
			slog.Error("postgres close failed", "error", err)
		}
	}()

	if err := postgres.Migrate(); err != nil {
		slog.Error("database migration error", "error", err)
		return
	}
	lg.Info("db connected")

	db := postgres.Get()

	skip, last, err := shouldSkip(db, time.Duration(*intervalSec)*time.Second)
	if err != nil {
		slog.Error("interval check failed", "error", err)
		return
	}
	if skip {
		lg.Info("skip: last run too recent",
			"finished_at", last.Format(time.RFC3339),
			"interval_seconds", *intervalSec,
		)
		return
	}
	lg.Info("interval check", "ok", true, "interval_seconds", *intervalSec)

	started := time.Now()
	stats, phase, err := runPipeline(ctx, db, lg, *upstreamDir, *protocolsJSON)
	if err != nil {
		recordFailure(db, lg, started, phase, err)
		return
	}

	row := newRefreshRun(started, stats, nil)
	if err := db.Create(row).Error; err != nil {
		slog.Error("refresh_run insert failed", "phase", "record-success", "error", err)
		return
	}

	lg.Info("refresh complete",
		"duration_ms", *row.DurationMs,
		"protocols", stats.Protocols,
		"adapter_files", stats.AdapterFiles,
		"adapter_files_skipped", stats.Skipped,
		"commits", 0,
	)
}

// runPipeline orchestrates walk + load + build.
// On failure it returns the phase name alongside the error
func runPipeline(ctx context.Context, db *gorm.DB, lg *slog.Logger, upstreamDir, protocolsJSON string) (dimensions.BuildStats, string, error) {
	adapters, err := dimensions.Walk(ctx, upstreamDir)
	if err != nil {
		return dimensions.BuildStats{}, "walk", err
	}
	lg.Info("walk done", "adapters", len(adapters))

	raw, err := dimensions.LoadProtocols(ctx, protocolsJSON)
	if err != nil {
		return dimensions.BuildStats{}, "load", err
	}
	lg.Info("load done", "data_files", len(raw))

	stats, err := dimensions.Build(ctx, db, raw, adapters, lg)
	if err != nil {
		return dimensions.BuildStats{}, "build", err
	}
	lg.Info("build done",
		"protocols", stats.Protocols,
		"adapter_files", stats.AdapterFiles,
		"skipped", stats.Skipped,
	)
	return stats, "", nil
}

// newRefreshRun builds a RefreshRun row with timestamps and duration filled in from started + time.Now().
func newRefreshRun(started time.Time, stats dimensions.BuildStats, errMsg *string) *models.RefreshRun {
	finished := time.Now()
	durMs := int(finished.Sub(started).Milliseconds())
	return &models.RefreshRun{
		StartedAt:           started,
		FinishedAt:          &finished,
		DurationMs:          &durMs,
		ProtocolsSeen:       stats.Protocols,
		AdapterFilesSeen:    stats.AdapterFiles,
		AdapterFilesSkipped: stats.Skipped,
		CommitsSeen:         0, // commits ingestion lands in a later step
		Error:               errMsg,
	}
}

// shouldSkip returns true when the most recent finished refresh_run is younger
// than interval. A zero interval disables the gate.
func shouldSkip(db *gorm.DB, interval time.Duration) (bool, time.Time, error) {
	if interval <= 0 {
		return false, time.Time{}, nil
	}
	var last models.RefreshRun
	err := db.Where("finished_at IS NOT NULL").
		Order("finished_at DESC").
		Limit(1).
		Take(&last).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, time.Time{}, nil
		}
		return false, time.Time{}, err
	}
	if last.FinishedAt == nil {
		return false, time.Time{}, nil
	}
	if time.Since(*last.FinishedAt) < interval {
		return true, *last.FinishedAt, nil
	}
	return false, time.Time{}, nil
}

// recordFailure writes a partial refresh_run row for an aborted pipeline so
// the gate in shouldSkip still sees forward progress and the operator can
// audit failures from the table.
func recordFailure(db *gorm.DB, lg *slog.Logger, started time.Time, phase string, cause error) {
	msg := fmt.Sprintf("%s: %v", phase, cause)
	row := newRefreshRun(started, dimensions.BuildStats{}, &msg)
	if err := db.Create(row).Error; err != nil {
		lg.Error("refresh_run failure insert failed",
			"phase", phase,
			"cause", cause,
			"insert_error", err,
		)
		return
	}
	lg.Error("refresh failed", "phase", phase, "error", cause, "duration_ms", *row.DurationMs)
}
