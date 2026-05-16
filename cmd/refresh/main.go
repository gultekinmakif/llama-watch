// One-shot orchestrator. Invoked by scripts/refresh.sh.
package main

import (
	"flag"
	"log"
	"log/slog"
	"time"

	"github.com/gultekinmakif/llama-watch/internal/config"
	"github.com/gultekinmakif/llama-watch/internal/db/postgres"
	"github.com/gultekinmakif/llama-watch/internal/dimensions"
	"github.com/gultekinmakif/llama-watch/internal/logger"
	"github.com/gultekinmakif/llama-watch/internal/models"
)

func main() {
	intervalSec := flag.Int("interval", 3300, "seconds; skip if the last refresh finished within this window")
	upstreamDir := flag.String("upstream-dir", "var/upstream", "parent of the cloned upstream repos")
	protocolsJSON := flag.String("protocols-json", "var/extracted/protocols.json", "path to the bun extractor output")
	snapshotOut := flag.String("snapshot-out", "var/snapshot/snapshot.json", "snapshot writer output path; accepted but not consumed in this step")
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
	lg.Debug("snapshot output path reserved for a later step", "snapshot_out", *snapshotOut)

	if err := postgres.New(cfg.DatabaseURL); err != nil {
		slog.Error("database connection error", "error", err)
		return
	}
	defer postgres.Close()

	if err := postgres.Migrate(); err != nil {
		slog.Error("database migration error", "error", err)
		return
	}
	lg.Info("db connected")

	db := postgres.Get()

	lg.Info("interval check", "ok", true, "interval_seconds", *intervalSec)

	started := time.Now()

	adapters, err := dimensions.Walk(*upstreamDir)
	if err != nil {
		return
	}
	lg.Info("walk done", "adapters", len(adapters))

	raw, err := dimensions.LoadProtocols(*protocolsJSON)
	if err != nil {
		return
	}
	lg.Info("load done", "data_files", len(raw))

	stats, err := dimensions.Build(db, raw, adapters, lg)
	if err != nil {
		return
	}
	lg.Info("build done",
		"protocols", stats.Protocols,
		"adapter_files", stats.AdapterFiles,
		"skipped", stats.Skipped,
	)

	finished := time.Now()
	durMs := int(finished.Sub(started).Milliseconds())
	row := models.RefreshRun{
		StartedAt:        started,
		FinishedAt:       &finished,
		DurationMs:       &durMs,
		ProtocolsSeen:    stats.Protocols,
		AdapterFilesSeen: stats.AdapterFiles,
		CommitsSeen:      0,
		Error:            nil,
	}
	if err := db.Create(&row).Error; err != nil {
		slog.Error("refresh_run insert failed", "error", err)
		return
	}

	lg.Info("refresh complete",
		"duration_ms", durMs,
		"protocols", stats.Protocols,
		"adapter_files", stats.AdapterFiles,
		"commits", 0,
	)
}
