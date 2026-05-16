// SQL query layer for /api/matrix: list protocols with limit/offset and a total count.
package api

import (
	"context"

	"github.com/gultekinmakif/llama-watch/internal/models"
	"github.com/gultekinmakif/llama-watch/internal/registry"
	"gorm.io/gorm"
)

// listProtocols returns the protocols page ordered by id for deterministic pagination.
// Each Row's Cells is zero-filled across the pinned column set and then flipped to 1
// for any dimension_kind present in adapter_files for that protocol.
func listProtocols(ctx context.Context, db *gorm.DB, limit, offset int) ([]Row, error) {
	var protos []models.Protocol
	if err := db.WithContext(ctx).
		Order("id").
		Limit(limit).
		Offset(offset).
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

// countProtocols returns the total number of protocols rows.
func countProtocols(ctx context.Context, db *gorm.DB) (int, error) {
	var n int64
	if err := db.WithContext(ctx).Model(&models.Protocol{}).Count(&n).Error; err != nil {
		return 0, err
	}
	return int(n), nil
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
