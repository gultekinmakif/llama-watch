// GET /api/dimensions handler. Returns the pinned dimension set with per-kind protocol coverage.
package api

import (
	"log/slog"
	"net/http"

	"github.com/gultekinmakif/llama-watch/internal/db/postgres"
	"github.com/gultekinmakif/llama-watch/internal/registry"
)

type DimensionEntry struct {
	Kind        string `json:"kind"`
	DisplayName string `json:"display_name"`
	Coverage    int    `json:"coverage"`
}

type dimensionsResponse struct {
	Dimensions []DimensionEntry `json:"dimensions"`
}

// Dimensions serves GET /api/dimensions.
func Dimensions(w http.ResponseWriter, r *http.Request) {
	db := postgres.Get()
	ctx := r.Context()

	coverage, err := dimensionCoverage(ctx, db)
	if err != nil {
		slog.Error("dimensions coverage failed", "error", err)
		writeErr(w, http.StatusInternalServerError, "internal", "internal error")
		return
	}

	cols := registry.Columns()
	out := make([]DimensionEntry, 0, len(cols))
	for _, c := range cols {
		out = append(out, DimensionEntry{
			Kind:        c.Key,
			DisplayName: c.Label,
			Coverage:    coverage[c.Key],
		})
	}

	writeJSON(w, http.StatusOK, dimensionsResponse{Dimensions: out})
}
