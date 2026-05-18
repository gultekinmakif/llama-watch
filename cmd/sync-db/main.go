// Invoked after tools/build-snapshot.ts.
// Loads var/snapshot/snapshot.json into the matrix, protocol_identities, and dim_file_coverage tables.
package main

import (
	"encoding/json"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/gultekinmakif/llama-watch/internal/config"
	"github.com/gultekinmakif/llama-watch/internal/db/postgres"
	"github.com/gultekinmakif/llama-watch/internal/logger"
	"github.com/gultekinmakif/llama-watch/internal/models"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type snapshotCell struct {
	Slug     string `json:"slug"`
	Metric   string `json:"metric"`
	CodePath string `json:"codePath"`
}

type snapshotProtocol struct {
	Slug     string   `json:"slug"`
	Name     string   `json:"name"`
	Category string   `json:"category"`
	Chains   []string `json:"chains"`
	DataFile string   `json:"dataFile"`
}

type snapshotFile struct {
	Cells     []snapshotCell     `json:"cells"`
	Protocols []snapshotProtocol `json:"protocols"`
}

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

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("getwd: %v", err)
	}
	snapshotPath := filepath.Join(cwd, "var", "snapshot", "snapshot.json")

	snap, err := readSnapshot(snapshotPath)
	if err != nil {
		log.Fatalf("read snapshot %s: %v", snapshotPath, err)
	}

	if err := postgres.New(cfg.DatabaseURL); err != nil {
		log.Fatalf("database connection: %v", err)
	}
	defer func() {
		if err := postgres.Close(); err != nil {
			slog.Error("postgres close failed", "error", err)
		}
	}()

	if err := postgres.Migrate(); err != nil {
		log.Fatalf("database migration: %v", err)
	}

	db := postgres.Get()
	identities := buildIdentities(snap.Protocols)
	rows := buildMatrix(snap.Cells)
	coverage := buildCoverage(snap.Cells)

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("TRUNCATE matrix, protocol_identities, dim_file_coverage").Error; err != nil {
			return err
		}
		if len(identities) > 0 {
			if err := tx.CreateInBatches(identities, 500).Error; err != nil {
				return err
			}
		}
		if len(rows) > 0 {
			if err := tx.CreateInBatches(rows, 500).Error; err != nil {
				return err
			}
		}
		if len(coverage) > 0 {
			if err := tx.CreateInBatches(coverage, 500).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("sync transaction: %v", err)
	}

	lg.Info("sync-db complete",
		"protocols", len(identities),
		"matrix_rows", len(rows),
		"coverage_rows", len(coverage),
		"snapshot", snapshotPath,
	)
}

func readSnapshot(path string) (snapshotFile, error) {
	var snap snapshotFile
	raw, err := os.ReadFile(path)
	if err != nil {
		return snap, err
	}
	if err := json.Unmarshal(raw, &snap); err != nil {
		return snap, err
	}
	return snap, nil
}

func buildIdentities(in []snapshotProtocol) []models.ProtocolIdentity {
	out := make([]models.ProtocolIdentity, len(in))
	for i, p := range in {
		out[i] = models.ProtocolIdentity{
			Slug:     p.Slug,
			Name:     p.Name,
			Category: nilIfEmpty(p.Category),
			Chains:   pq.StringArray(p.Chains),
			DataFile: nilIfEmpty(p.DataFile),
		}
	}
	return out
}

// Empty strings in the snapshot represent "not set" upstream; persist as SQL NULL
// rather than "" so IS NULL queries work and the column's NULL semantics survive.
func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func buildMatrix(in []snapshotCell) []models.Matrix {
	out := make([]models.Matrix, len(in))
	for i, c := range in {
		out[i] = models.Matrix{
			Slug:     c.Slug,
			Metric:   c.Metric,
			CodePath: c.CodePath,
		}
	}
	return out
}

func buildCoverage(in []snapshotCell) []models.DimFileCoverage {
	seen := make(map[string]map[string]struct{}, len(in))
	for _, c := range in {
		if c.CodePath == "" {
			continue
		}
		set, ok := seen[c.CodePath]
		if !ok {
			set = make(map[string]struct{})
			seen[c.CodePath] = set
		}
		set[c.Metric] = struct{}{}
	}
	out := make([]models.DimFileCoverage, 0)
	for codePath, metrics := range seen {
		for metric := range metrics {
			out = append(out, models.DimFileCoverage{CodePath: codePath, Metric: metric})
		}
	}
	return out
}
