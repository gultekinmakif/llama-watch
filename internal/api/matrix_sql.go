// SQL query layer for /api/matrix.
package api

import (
	"context"

	"github.com/gultekinmakif/llama-watch/internal/models"
	"github.com/gultekinmakif/llama-watch/internal/registry"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// listProtocols returns the protocol_identities page ordered by q.Sort then slug.
// Each Row's Cells is zero-filled across the pinned column set and then flipped to 1
// for any metric present in matrix for that slug.
func listProtocols(ctx context.Context, db *gorm.DB, q MatrixQuery) ([]Row, error) {
	tx := db.WithContext(ctx).Model(&models.ProtocolIdentity{})
	tx = applyMatrixFilters(tx, q)
	tx = applyMatrixOrder(tx, q)

	var idents []models.ProtocolIdentity
	if err := tx.
		Limit(q.Limit).
		Offset(q.Offset).
		Find(&idents).Error; err != nil {
		return nil, err
	}

	slugs := make([]string, len(idents))
	for i, p := range idents {
		slugs[i] = p.Slug
	}

	cellsBySlug, err := fetchCells(ctx, db, slugs)
	if err != nil {
		return nil, err
	}

	rows := make([]Row, 0, len(idents))
	for _, p := range idents {
		cells := initCells()
		for metric := range cellsBySlug[p.Slug] {
			if _, ok := cells[metric]; ok {
				cells[metric] = 1
			}
		}
		rows = append(rows, Row{
			Slug:     p.Slug,
			Name:     p.Name,
			Category: p.Category,
			Chains:   []string(p.Chains),
			Cells:    cells,
		})
	}
	return rows, nil
}

// countProtocols returns the total number of protocol_identities rows after filtering.
func countProtocols(ctx context.Context, db *gorm.DB, q MatrixQuery) (int, error) {
	tx := db.WithContext(ctx).Model(&models.ProtocolIdentity{})
	tx = applyMatrixFilters(tx, q)

	var n int64
	if err := tx.Count(&n).Error; err != nil {
		return 0, err
	}
	return int(n), nil
}

// applyMatrixOrder layers the requested sort plus a slug tiebreaker onto tx.
// NULL category sorts last regardless of direction.
func applyMatrixOrder(tx *gorm.DB, q MatrixQuery) *gorm.DB {
	switch q.Sort {
	case "name":
		tx = tx.Order("name " + q.Order)
	case "category":
		tx = tx.Order("category " + q.Order + " NULLS LAST")
	case "coverage":
		tx = tx.Order("(SELECT COUNT(*) FROM matrix WHERE matrix.slug = protocol_identities.slug) " + q.Order)
	}
	return tx.Order("slug")
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

// fetchCells returns slug -> {metric: 1} for the given slugs.
// Only present metrics are included; the caller zero-fills the closed columns set.
func fetchCells(ctx context.Context, db *gorm.DB, slugs []string) (map[string]map[string]int, error) {
	if len(slugs) == 0 {
		return map[string]map[string]int{}, nil
	}
	var rows []struct {
		Slug   string `gorm:"column:slug"`
		Metric string `gorm:"column:metric"`
	}
	if err := db.WithContext(ctx).
		Table("matrix").
		Select("slug, metric").
		Where("slug IN ?", slugs).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make(map[string]map[string]int, len(slugs))
	for _, r := range rows {
		m, ok := out[r.Slug]
		if !ok {
			m = make(map[string]int)
			out[r.Slug] = m
		}
		m[r.Metric] = 1
	}
	return out, nil
}

// initCells returns a fresh cells map with every pinned column key set to 0.
func initCells() map[string]int {
	cols := registry.Columns()
	cells := make(map[string]int, len(cols))
	for _, c := range cols {
		cells[c.Key] = 0
	}
	return cells
}
