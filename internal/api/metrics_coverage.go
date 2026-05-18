// GET /api/metrics-coverage handler. Returns per-adapter-file emission rows projected from matrix.
package api

import (
	"net/http"
)

type MetricsCoverageEntry struct {
	CodePath string `json:"code_path"`
	Metric   string `json:"metric"`
}

type metricsCoverageResponse struct {
	Rows  []MetricsCoverageEntry `json:"rows"`
	Total int                    `json:"total"`
}

// MetricsCoverage serves GET /api/metrics-coverage.
func MetricsCoverage(w http.ResponseWriter, r *http.Request) {

}
