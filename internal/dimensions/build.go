// Pass 0 of the slug-join algorithm.
// Reads the extracted protocols JSON plus the walker output, and writes
// protocols and adapter_files rows inside one transaction.
package dimensions

import (
	"log/slog"
	"path"
	"strings"

	"github.com/gultekinmakif/llama-watch/internal/models"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BuildStats struct {
	Protocols    int
	AdapterFiles int
	// Skipped also counts adapter files referenced by data{N}.ts but missing on disk.
	Skipped int
}

// Build runs Pass 1 of the slug-join algorithm in one DB transaction.
// raw is the LoadProtocols output keyed by source data file ("data1", ..., "data6").
// adapters is the Walk output of every adapter file on disk.
// log may be nil; the default slog handler is used in that case.
func Build(db *gorm.DB, raw map[string][]RawProtocol, adapters []Adapter, log *slog.Logger) (BuildStats, error) {
	if log == nil {
		log = slog.Default()
	}

	byRel := indexAdapters(adapters)

	var dims map[string]uint64
	if err := db.Transaction(func(tx *gorm.DB) error {
		var err error
		dims, err = loadDimensionIDs(tx)
		return err
	}); err != nil {
		return BuildStats{}, err
	}

	var stats BuildStats
	err := db.Transaction(func(tx *gorm.DB) error {
		for src, list := range raw {
			dataFile := src + ".ts"
			for _, rp := range list {
				if err := buildOne(tx, log, byRel, dims, rp, dataFile, &stats); err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return BuildStats{}, err
	}
	return stats, nil
}

// indexAdapters builds an O(1) lookup keyed by RelPath.
func indexAdapters(adapters []Adapter) map[string]Adapter {
	out := make(map[string]Adapter, len(adapters))
	for _, a := range adapters {
		out[a.RelPath] = a
	}
	return out
}

// buildOne handles a single RawProtocol: upsert the protocol row and resolve
// its TVL adapter plus every dimension-adapter file it references.
func buildOne(
	tx *gorm.DB,
	log *slog.Logger,
	byRel map[string]Adapter,
	dims map[string]uint64,
	rp RawProtocol,
	dataFile string,
	stats *BuildStats,
) error {
	p, err := upsertProtocol(tx, rp, dataFile)
	if err != nil {
		return err
	}
	stats.Protocols++

	if err := resolveTVL(tx, log, byRel, p.ID, rp.Module, stats); err != nil {
		return err
	}

	for dimType, dimSlug := range rp.Dimensions {
		if err := resolveDimension(tx, log, byRel, dims, p.ID, dimType, dimSlug, stats); err != nil {
			return err
		}
	}
	return nil
}

// upsertProtocol inserts or updates a protocols row keyed by canonical slug.
// Updates on conflict cover name, category, chains, and data_file so a later
// data{N}.ts file overrides an earlier one deterministically.
func upsertProtocol(tx *gorm.DB, rp RawProtocol, dataFile string) (*models.Protocol, error) {
	chains := normalizeChains(rp.Chains)
	df := dataFile
	var category *string
	if rp.Category != "" {
		c := rp.Category
		category = &c
	}

	p := models.Protocol{
		Slug:     Canonical(rp.Name),
		Name:     rp.Name,
		Category: category,
		Chains:   chains,
		DataFile: &df,
	}

	err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "slug"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "category", "chains", "data_file"}),
	}).Create(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// resolveTVL looks up the DefiLlama-Adapters file for a protocol module.
// Missing files are logged and counted as skipped; no row is created.
func resolveTVL(
	tx *gorm.DB,
	log *slog.Logger,
	byRel map[string]Adapter,
	protocolID uint64,
	module string,
	stats *BuildStats,
) error {
	rel := "DefiLlama-Adapters/projects/" + module
	a, ok := byRel[rel]
	if !ok {
		log.Warn("tvl adapter missing on disk", "module", module)
		stats.Skipped++
		return nil
	}
	pid := protocolID
	row := models.AdapterFile{
		ProtocolID:    &pid,
		DimensionID:   nil,
		Repo:          "defillama-adapters",
		Path:          "projects/" + module,
		Slug:          moduleStem(module),
		DimensionKind: "tvl",
		Orphan:        false,
	}
	if err := upsertAdapterFile(tx, &row); err != nil {
		return err
	}
	stats.AdapterFiles++
	_ = a
	return nil
}

// loadDimensionIDs reads the dimensions table into kind -> id so adapter_files
// rows can be populated without a per-row query. The dimensions seed is written
// by postgres.Migrate before Build ever runs.
func loadDimensionIDs(tx *gorm.DB) (map[string]uint64, error) {
	var rows []models.Dimension
	if err := tx.Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make(map[string]uint64, len(rows))
	for _, r := range rows {
		out[r.Kind] = r.ID
	}
	return out, nil
}

// resolveDimension locates a dimension-adapter file for a (type, slug) pair,
// detects which sub-metric kinds it emits, and inserts one row per kind.
func resolveDimension(
	tx *gorm.DB,
	log *slog.Logger,
	byRel map[string]Adapter,
	dims map[string]uint64,
	protocolID uint64,
	dimType, dimSlug string,
	stats *BuildStats,
) error {
	candidates := []string{
		"dimension-adapters/" + dimType + "/" + dimSlug + ".ts",
		"dimension-adapters/" + dimType + "/" + dimSlug + "/index.ts",
	}
	var resolved Adapter
	found := false
	for _, c := range candidates {
		if a, ok := byRel[c]; ok {
			resolved = a
			found = true
			break
		}
	}
	if !found {
		log.Warn("dimension adapter missing on disk", "type", dimType, "slug", dimSlug)
		stats.Skipped++
		return nil
	}

	kinds, err := DetectKinds(resolved.AbsPath, dimType)
	if err != nil {
		return err
	}
	if len(kinds) == 0 {
		return nil
	}

	relPath := strings.TrimPrefix(resolved.RelPath, "dimension-adapters/")
	for _, kind := range kinds {
		did, ok := dims[kind]
		if !ok {
			// Dimension seed is missing this kind. Surface it and skip the row
			// rather than violate the FK.
			log.Warn("dimension kind not seeded", "kind", kind)
			continue
		}
		pid := protocolID
		dimID := did
		row := models.AdapterFile{
			ProtocolID:    &pid,
			DimensionID:   &dimID,
			Repo:          "dimension-adapters",
			Path:          relPath,
			Slug:          resolved.Slug,
			DimensionKind: kind,
			Orphan:        false,
		}
		if err := upsertAdapterFile(tx, &row); err != nil {
			return err
		}
		stats.AdapterFiles++
	}
	return nil
}

// upsertAdapterFile uses ON CONFLICT against the (repo, path, dimension_kind)
// unique index so re-runs are idempotent and the protocol_id / dimension_id
// associations stay current.
func upsertAdapterFile(tx *gorm.DB, row *models.AdapterFile) error {
	return tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "repo"},
			{Name: "path"},
			{Name: "dimension_kind"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"protocol_id",
			"dimension_id",
			"slug",
			"orphan",
		}),
	}).Create(row).Error
}

// normalizeChains lowercases every entry and returns a pq.StringArray, which is
// the column type the protocols.chains field expects.
func normalizeChains(chains []string) pq.StringArray {
	out := make(pq.StringArray, len(chains))
	for i, c := range chains {
		out[i] = strings.ToLower(c)
	}
	return out
}

// moduleStem extracts the canonical stem from a module reference like
// "uniswap-v2/index.js" or "wbtc.js". It is intentionally lighter than
// FromFilename: no separator normalization, just basename + extension strip.
func moduleStem(module string) string {
	s := strings.TrimRight(module, "/")
	if s == "" {
		return ""
	}
	base := path.Base(s)
	base = strings.TrimSuffix(base, ".js")
	base = strings.TrimSuffix(base, ".ts")
	if base == "index" {
		parent := path.Dir(s)
		if parent != "." && parent != "/" {
			return path.Base(parent)
		}
	}
	return base
}
