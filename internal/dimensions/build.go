// Pass 1 of the slug-join algorithm.
package dimensions

import (
	"strings"

	"github.com/gultekinmakif/llama-watch/internal/models"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Build runs Pass 1 of the slug-join algorithm in one DB transaction.
func Build(db *gorm.DB, raw map[string][]RawProtocol, adapters []Adapter) error {
	byRel := indexAdapters(adapters)

	err := db.Transaction(func(tx *gorm.DB) error {
		for src, list := range raw {
			dataFile := src + ".ts"
			for _, rp := range list {
				if err := buildOne(tx, byRel, rp, dataFile); err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// indexAdapters builds an O(1) lookup keyed by RelPath.
func indexAdapters(adapters []Adapter) map[string]Adapter {
	out := make(map[string]Adapter, len(adapters))
	for _, a := range adapters {
		out[a.RelPath] = a
	}
	return out
}

// buildOne handles a single RawProtocol
func buildOne(
	tx *gorm.DB,
	byRel map[string]Adapter,
	rp RawProtocol,
	dataFile string,
) error {
	p, err := insertProtocol(tx, rp, dataFile)
	if err != nil {
		return err
	}

	if err := resolveTVL(tx, byRel, p.ID, rp.Module); err != nil {
		return err
	}

	return nil
}

// insertProtocol inserts or updates a protocols row keyed by canonical slug.
func insertProtocol(tx *gorm.DB, rp RawProtocol, dataFile string) (*models.Protocol, error) {
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
func resolveTVL(
	tx *gorm.DB,
	byRel map[string]Adapter,
	protocolID uint64,
	module string,
) error {
	rel := "DefiLlama-Adapters/projects/" + module
	a, ok := byRel[rel]
	if !ok {
		return nil
	}
	pid := protocolID
	row := models.AdapterFile{
		ProtocolID:    &pid,
		DimensionID:   nil,
		Repo:          "defillama-adapters",
		Path:          "projects/" + module,
		Slug:          module,
		DimensionKind: "tvl",
		Orphan:        false,
	}
	if err := insertAdapterFile(tx, &row); err != nil {
		return err
	}
	_ = a
	return nil
}

// insertAdapterFile uses ON CONFLICT against the (repo, path, dimension_kind)
// unique index so re-runs are idempotent and the protocol_id / dimension_id
// associations stay current.
func insertAdapterFile(tx *gorm.DB, row *models.AdapterFile) error {
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

// normalizeChains lowercases every entry and returns a pq.StringArray, protocols.chains expects.
func normalizeChains(chains []string) pq.StringArray {
	out := make(pq.StringArray, len(chains))
	for i, c := range chains {
		out[i] = strings.ToLower(c)
	}
	return out
}
