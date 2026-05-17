// SQL query layer for /api/dimensions.
package api

import (
	"context"

	"gorm.io/gorm"
)

const dimensionCoverageSQL = `
SELECT dimension_kind AS kind, COUNT(DISTINCT protocol_id) AS coverage
FROM adapter_files
WHERE orphan = false
GROUP BY dimension_kind
`

// dimensionCoverage returns COUNT(DISTINCT protocol_id) per dimension_kind, excluding orphans.
// Result is keyed by dimension_kind; kinds with no matching rows are absent (caller treats absence as zero).
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
