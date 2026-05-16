// 1.4: Pass 0 of the slug-join algorithm.
// Third pipeline phase. Joins the extracted protocols JSON with the walker output
// and writes protocols and adapter_files rows inside one transaction.
package dimensions

import (
	"context"
	"log/slog"
	"strings"

	"github.com/gultekinmakif/llama-watch/internal/models"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	tvlAdapterPrefix       = "DefiLlama-Adapters/projects/"
	dimensionAdapterPrefix = "dimension-adapters/"
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
// ctx is bound to the gorm session so cancellation propagates into every query
// and is also re-checked at the top of each per-protocol iteration.
func Build(ctx context.Context, db *gorm.DB, raw map[string][]RawProtocol, adapters []Adapter, log *slog.Logger) (BuildStats, error) {
	if log == nil {
		log = slog.Default()
	}

	byRel := indexAdapters(adapters)
	cdb := db.WithContext(ctx)

	var stats BuildStats
	err := cdb.Transaction(func(tx *gorm.DB) error {
		for src, list := range raw {
			dataFile := src + ".ts"
			for _, rp := range list {
				if cerr := ctx.Err(); cerr != nil {
					return cerr
				}
				if err := buildOne(tx, log, byRel, rp, dataFile, &stats); err != nil {
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
		if err := resolveDimension(tx, log, byRel, p.ID, dimType, dimSlug, stats); err != nil {
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
	rel := tvlAdapterPrefix + module
	_, ok := byRel[rel]
	if !ok {
		log.Warn("tvl adapter missing on disk", "module", module)
		stats.Skipped++
		return nil
	}
	pid := protocolID
	row := models.AdapterFile{
		ProtocolID:    &pid,
		Repo:          "defillama-adapters",
		Path:          "projects/" + module,
		Slug:          pathSlug(module),
		DimensionKind: "tvl",
		Orphan:        false,
	}
	if err := upsertAdapterFile(tx, &row); err != nil {
		return err
	}
	stats.AdapterFiles++
	return nil
}

// dimensionCandidates returns the relative paths the walker would emit for
// a given dimension type and slug, in resolution order.
//
// > The first existing path in the walker output wins!
func dimensionCandidates(dimType, dimSlug string) []string {
	base := dimensionAdapterPrefix + dimType + "/" + dimSlug
	return []string{base + ".ts", base + "/index.ts"}
}

// findDimensionAdapter returns the first Adapter whose RelPath matches one
// of the candidate paths for (dimType, dimSlug), or nil.
func findDimensionAdapter(byRel map[string]Adapter, dimType, dimSlug string) *Adapter {
	for _, candidate := range dimensionCandidates(dimType, dimSlug) {
		if a, ok := byRel[candidate]; ok {
			return &a
		}
	}
	return nil
}

// resolveDimension locates a dimension-adapter file for a (type, slug) pair,
// detects which sub-metric kinds it emits, and inserts one row per kind.
func resolveDimension(
	tx *gorm.DB,
	log *slog.Logger,
	byRel map[string]Adapter,
	protocolID uint64,
	dimType, dimSlug string,
	stats *BuildStats,
) error {
	resolved := findDimensionAdapter(byRel, dimType, dimSlug)
	if resolved == nil {
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

	relPath := strings.TrimPrefix(resolved.RelPath, dimensionAdapterPrefix)
	for _, kind := range kinds {
		pid := protocolID
		row := models.AdapterFile{
			ProtocolID:    &pid,
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
// unique index so re-runs are idempotent and the protocol_id association stays current.
func upsertAdapterFile(tx *gorm.DB, row *models.AdapterFile) error {
	return tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "repo"},
			{Name: "path"},
			{Name: "dimension_kind"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"protocol_id",
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
