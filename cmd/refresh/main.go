// One-shot orchestrator. Invoked by scripts/refresh.sh.
package main

import (
	"flag"
	"log"
	"log/slog"

	"github.com/gultekinmakif/llama-watch/internal/config"
	"github.com/gultekinmakif/llama-watch/internal/db/postgres"
	"github.com/gultekinmakif/llama-watch/internal/dimensions"
	"github.com/gultekinmakif/llama-watch/internal/logger"
	"github.com/gultekinmakif/llama-watch/internal/models"
)

func main() {
	upstreamDir := "var/upstream"
	protocolsJSON := "var/extracted/protocols.json"
	flag.Parse()

	cfg, err := config.Load()

	if err := postgres.New(*cfg.DatabaseURL); err != nil {
		return
	}
	defer postgres.Close()

	if err := postgres.Migrate(); err != nil {
		return
	}

	db := postgres.Get()

	adapters, err := dimensions.Walk(upstreamDir)
	if err != nil {
		return
	}

	raw, err := dimensions.LoadProtocols(protocolsJSON)
	if err != nil {
		return
	}

	lg, err := logger.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	slog.SetDefault(lg)
	stats, err := dimensions.Build(db, raw, adapters, lg)
	if err != nil {
		return
	}

	row := models.RefreshRun{
		ProtocolsSeen:    stats.Protocols,
		AdapterFilesSeen: stats.AdapterFiles,
		CommitsSeen:      0,
		Error:            nil,
	}
	if err := db.Create(&row).Error; err != nil {
		return
	}

}
