// SQL query layer for /api/metrics-coverage.
package api

import (
	"context"
	"strings"

	"gorm.io/gorm"
)

// listMetricsCoverage returns distinct (code_path, metric) rows from matrix,
// optionally filtered by metric and by protocol slug. Empty filters are ignored.
func listMetricsCoverage(ctx context.Context, db *gorm.DB, metric, protocol string) ([]MetricsCoverageEntry, error) {
	sql := "SELECT DISTINCT code_path, metric FROM matrix WHERE " + strings.Join(clauses, " AND ") + " ORDER BY code_path, metric"

	var args []any
	clauses := []string{"code_path != ''"}
	if metric != "" {
		args = append(args, metric)
		clauses = append(clauses, "metric = ?")
	}
	if protocol != "" {
		args = append(args, protocol)
		clauses = append(clauses, "slug = ?")
	}

	rows := make([]MetricsCoverageEntry, 0)
	if err := db.WithContext(ctx).Raw(sql, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return nil, nil
}
