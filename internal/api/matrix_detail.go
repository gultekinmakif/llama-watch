// GET /api/matrix/{slug} handler. Returns identity, chains, and per-dimension coverage.
package api

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gultekinmakif/llama-watch/internal/db/postgres"
	"gorm.io/gorm"
)

// MatrixDetail serves GET /api/matrix/{slug}.
func MatrixDetail(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	db := postgres.Get()
	ctx := r.Context()

	detail, err := fetchMatrixDetail(ctx, db, slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeErr(w, http.StatusNotFound, "not_found", "protocol not found")
			return
		}
		slog.Error("matrix detail fetch failed", "slug", slug, "error", err)
		writeErr(w, http.StatusInternalServerError, "internal", "internal error")
		return
	}

	writeJSON(w, http.StatusOK, detail)
}

type ProtocolDimension struct {
	Kind      string  `json:"kind"`
	Present   bool    `json:"present"`
	GitHubURL *string `json:"github_url,omitempty"`
}

type ProtocolDetail struct {
	Slug       string              `json:"slug"`
	Name       string              `json:"name"`
	Category   *string             `json:"category,omitempty"`
	Chains     []string            `json:"chains"`
	Dimensions []ProtocolDimension `json:"dimensions"`
}
