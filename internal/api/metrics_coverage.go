// GET /api/metrics-coverage handler. Returns per-adapter-file emission rows projected from matrix.
package api

import (
	"log/slog"
	"net/http"

	"github.com/gultekinmakif/llama-watch/internal/db/postgres"
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
	db := postgres.Get()
	ctx := r.Context()

	q := r.URL.Query()
	metric := q.Get("metric")
	protocol := q.Get("protocol")

	rows, err := listMetricsCoverage(ctx, db, metric, protocol)
	if err != nil {
		slog.Error("metrics coverage failed", "error", err)
		writeErr(w, http.StatusInternalServerError, "internal", "internal error")
		return
	}

	writeJSON(w, http.StatusOK, metricsCoverageResponse{Rows: rows, Total: len(rows)})
}
