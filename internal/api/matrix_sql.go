// SQL query layer for /api/matrix: list protocols with limit/offset and a total count.
package api

import (
	"context"

	"github.com/gultekinmakif/llama-watch/internal/models"
	"gorm.io/gorm"
)

// listProtocols returns the protocols page ordered by id for deterministic pagination.
// Each Row's Cells is set to a non-nil empty map; cells assembly is a later step.
func listProtocols(ctx context.Context, db *gorm.DB, limit, offset int) ([]Row, error) {
	var protos []models.Protocol
	if err := db.WithContext(ctx).
		Order("id").
		Limit(limit).
		Offset(offset).
		Find(&protos).Error; err != nil {
		return nil, err
	}

	rows := make([]Row, 0, len(protos))
	for _, p := range protos {
		rows = append(rows, Row{
			Slug:     p.Slug,
			Name:     p.Name,
			Category: p.Category,
			TVL:      p.TVL,
			Chains:   []string(p.Chains),
			Cells:    map[string]int{},
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
