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
	clauses := []string{"code_path != ''"}
	var args []any
	if metric != "" {
		clauses = append(clauses, "metric = ?")
		args = append(args, metric)
	}
	if protocol != "" {
		clauses = append(clauses, "slug = ?")
		args = append(args, protocol)
	}
	sql := "SELECT DISTINCT code_path, metric FROM matrix WHERE " + strings.Join(clauses, " AND ") + " ORDER BY code_path, metric"

	rows := make([]MetricsCoverageEntry, 0)
	if err := db.WithContext(ctx).Raw(sql, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}
