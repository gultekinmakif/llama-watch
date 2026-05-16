// SQL query layer for /api/matrix.
package api

import (
	"context"

	"github.com/gultekinmakif/llama-watch/internal/models"
	"github.com/gultekinmakif/llama-watch/internal/registry"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// listProtocols returns the protocols page ordered by q.Sort then id for deterministic pagination.
// Each Row's Cells is zero-filled across the pinned column set and then flipped to 1
// for any dimension_kind present in adapter_files for that protocol.
func listProtocols(ctx context.Context, db *gorm.DB, q MatrixQuery) ([]Row, error) {
	tx := db.WithContext(ctx).Model(&models.Protocol{})
	tx = applyMatrixFilters(tx, q)
	tx = applyMatrixOrder(tx, q)

	var protos []models.Protocol
	if err := tx.
		Limit(q.Limit).
		Offset(q.Offset).
		Find(&protos).Error; err != nil {
		return nil, err
	}

	protoIDs := make([]uint64, len(protos))
	for i, p := range protos {
		protoIDs[i] = p.ID
	}

	cellsByID, err := fetchCells(ctx, db, protoIDs)
	if err != nil {
		return nil, err
	}

	rows := make([]Row, 0, len(protos))
	for _, p := range protos {
		cells := initCells()
		for kind := range cellsByID[p.ID] {
			if _, ok := cells[kind]; ok {
				cells[kind] = 1
			}
		}
		rows = append(rows, Row{
			Slug:     p.Slug,
			Name:     p.Name,
			Category: p.Category,
			TVL:      p.TVL,
			Chains:   []string(p.Chains),
			Cells:    cells,
		})
	}
	return rows, nil
}

// countProtocols returns the total number of protocols rows after filtering.
// Uses the same filter chain as listProtocols so total tracks the visible set.
func countProtocols(ctx context.Context, db *gorm.DB, q MatrixQuery) (int, error) {
	tx := db.WithContext(ctx).Model(&models.Protocol{})
	tx = applyMatrixFilters(tx, q)

	var n int64
	if err := tx.Count(&n).Error; err != nil {
		return 0, err
	}
	return int(n), nil
}

// applyMatrixOrder layers the requested sort plus an id tiebreaker onto tx.
// NULL tvl and NULL category sort last regardless of direction so they cluster.
// Coverage is computed from a correlated COUNT against adapter_files filtered to non-orphan rows.
// Empty q.Sort (snapshot path) falls through the switch and yields the id-only ordering.
func applyMatrixOrder(tx *gorm.DB, q MatrixQuery) *gorm.DB {
	switch q.Sort {
	case "name":
		tx = tx.Order("name " + q.Order)
	case "category":
		tx = tx.Order("category " + q.Order + " NULLS LAST")
	case "tvl":
		tx = tx.Order("tvl " + q.Order + " NULLS LAST")
	case "coverage":
		tx = tx.Order("(SELECT COUNT(*) FROM adapter_files WHERE adapter_files.protocol_id = protocols.id AND adapter_files.orphan = false) " + q.Order)
	}
	return tx.Order("id")
}

// applyMatrixFilters layers the chains / categories / q WHERE clauses onto tx.
// Empty slices and empty strings short-circuit each branch so a zero-value
// MatrixQuery is equivalent to "no filtering" for the snapshot path.
func applyMatrixFilters(tx *gorm.DB, q MatrixQuery) *gorm.DB {
	if len(q.Chains) > 0 {
		tx = tx.Where("chains && ?", pq.StringArray(q.Chains))
	}
	if len(q.Categories) > 0 {
		tx = tx.Where("category IN ?", q.Categories)
	}
	if q.Q != "" {
		pattern := "%" + q.Q + "%"
		tx = tx.Where("slug ILIKE ? OR name ILIKE ?", pattern, pattern)
	}
	return tx
}

// fetchCells returns protocol_id -> {dimension_kind: 1} for the given protocol IDs.
// Only present kinds are included; the caller zero-fills the closed columns set.
func fetchCells(ctx context.Context, db *gorm.DB, protoIDs []uint64) (map[uint64]map[string]int, error) {
	byID, err := fetchAdapterRows(ctx, db, protoIDs)
	if err != nil {
		return nil, err
	}
	out := make(map[uint64]map[string]int, len(byID))
	for pid, rows := range byID {
		m := make(map[string]int, len(rows))
		for _, r := range rows {
			m[r.DimensionKind] = 1
		}
		out[pid] = m
	}
	return out, nil
}

// initCells returns a fresh cells map with every pinned column key set to 0.
// Sources the closed set from the registry so the result tracks one source of truth.
func initCells() map[string]int {
	cols := registry.Columns()
	cells := make(map[string]int, len(cols))
	for _, c := range cols {
		cells[c.Key] = 0
	}
	return cells
}
