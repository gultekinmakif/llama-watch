// SQL query layer for /api/metrics-coverage.
package api

import (
	"context"

	"gorm.io/gorm"
)

// listMetricsCoverage returns distinct (code_path, metric) rows from matrix,
// optionally filtered by metric and by protocol slug. Empty filters are ignored.
func listMetricsCoverage(ctx context.Context, db *gorm.DB, metric, protocol string) ([]MetricsCoverageEntry, error) {
	return nil, nil
}
