// Invoked after tools/build-snapshot.ts.
// Loads var/snapshot/snapshot.json into the matrix and protocol_identities tables.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

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

type tvlProtocol struct {
	Module string  `json:"module"`
	TVL    float64 `json:"tvl"`
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
	// SNAPSHOT_OUT mirrors the env contract scripts/refresh.sh and tools/build-snapshot.ts honor.
	snapshotPath := os.Getenv("SNAPSHOT_OUT")
	if snapshotPath == "" {
		snapshotPath = "var/snapshot/snapshot.json"
	}
	if !filepath.IsAbs(snapshotPath) {
		snapshotPath = filepath.Join(cwd, snapshotPath)
	}

	snap, err := readSnapshot(snapshotPath)
	if err != nil {
		log.Fatalf("read snapshot %s: %v", snapshotPath, err)
	}

	// Catch duplicate slugs before TRUNCATE so a bad snapshot does not leave the
	// matrix table half-empty if CreateInBatches trips a primary-key conflict.
	if err := checkDuplicateSlugs(snap.Protocols); err != nil {
		log.Fatalf("snapshot validation: %v", err)
	}

	tvlBySlug, err := loadTVL(cwd)
	if err != nil {
		log.Fatalf("load tvl.json: %v", err)
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
	identities := buildIdentities(snap.Protocols, tvlBySlug)
	rows := buildMatrix(snap.Cells)

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("TRUNCATE matrix, protocol_identities").Error; err != nil {
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
		return nil
	})
	if err != nil {
		log.Fatalf("sync transaction: %v", err)
	}

	tvlCount := 0
	for _, id := range identities {
		if id.TvlUSD != nil {
			tvlCount++
		}
	}
	lg.Info("sync-db complete",
		"protocols", len(identities),
		"protocols_with_tvl", tvlCount,
		"matrix_rows", len(rows),
		"snapshot", snapshotPath,
	)
}

func checkDuplicateSlugs(protocols []snapshotProtocol) error {
	seen := make(map[string]struct{}, len(protocols))
	for _, p := range protocols {
		if _, ok := seen[p.Slug]; ok {
			return fmt.Errorf("duplicate slug %q in snapshot.protocols", p.Slug)
		}
		seen[p.Slug] = struct{}{}
	}
	return nil
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

func buildIdentities(in []snapshotProtocol, tvl map[string]float64) []models.ProtocolIdentity {
	out := make([]models.ProtocolIdentity, len(in))
	for i, p := range in {
		out[i] = models.ProtocolIdentity{
			Slug:     p.Slug,
			Name:     p.Name,
			Category: nilIfEmpty(p.Category),
			Chains:   pq.StringArray(p.Chains),
			DataFile: nilIfEmpty(p.DataFile),
		}
		if tvl != nil {
			if v, ok := tvl[p.Slug]; ok {
				out[i].TvlUSD = &v
			}
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

func loadTVL(cwd string) (map[string]float64, error) {
	tvlPath := os.Getenv("TVL_PATH")
	if tvlPath == "" {
		tvlPath = "var/snapshot/tvl.json"
	}
	if !filepath.IsAbs(tvlPath) {
		tvlPath = filepath.Join(cwd, tvlPath)
	}
	raw, err := os.ReadFile(tvlPath)
	if err != nil {
		if os.IsNotExist(err) {
			slog.Warn("tvl.json not found, skipping TVL population", "path", tvlPath)
			return nil, nil
		}
		return nil, err
	}
	var protocols []tvlProtocol
	if err := json.Unmarshal(raw, &protocols); err != nil {
		return nil, fmt.Errorf("parse tvl.json: %w", err)
	}
	out := make(map[string]float64, len(protocols))
	for _, p := range protocols {
		slug := normalizeModule(p.Module)
		if slug != "" && p.TVL > 0 {
			out[slug] = p.TVL
		}
	}
	return out, nil
}

// Mirrors normalizeSlug in tools/build-snapshot.ts; both must stay in sync.
func normalizeModule(mod string) string {
	s := mod
	if strings.HasSuffix(s, ".js") {
		s = s[:len(s)-3]
	} else if strings.HasSuffix(s, ".ts") {
		s = s[:len(s)-3]
	}
	if strings.HasSuffix(s, "/index") {
		s = s[:len(s)-6]
	}
	return s
}
