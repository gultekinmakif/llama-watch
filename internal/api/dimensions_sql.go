// SQL query layer for /api/dimensions.
package api

import (
	"context"

	"gorm.io/gorm"
)

const dimensionCoverageSQL = `
SELECT metric AS kind, COUNT(*) AS coverage
FROM matrix
GROUP BY metric
`

// dimensionCoverage returns COUNT(*) per metric from the matrix table.
// Result is keyed by metric; metrics with no rows are absent (caller treats absence as zero).
func dimensionCoverage(ctx context.Context, db *gorm.DB) (map[string]int, error) {
	var rows []struct {
		Kind     string
		Coverage int
	}
	if err := db.WithContext(ctx).Raw(dimensionCoverageSQL).Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := make(map[string]int, len(rows))
	for _, r := range rows {
		out[r.Kind] = r.Coverage
	}
	return out, nil
}
