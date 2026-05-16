// GET /api/matrix/{slug} handler. Returns identity, chains, and per-dimension coverage.
package api

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

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

type LastCommit struct {
	SHA         string    `json:"sha"`
	Author      string    `json:"author"`
	CommittedAt time.Time `json:"committed_at"`
	GitHubURL   string    `json:"github_url"`
}

type ProtocolDimension struct {
	Kind       string      `json:"kind"`
	Present    bool        `json:"present"`
	FilePath   *string     `json:"file_path"`
	Repo       *string     `json:"repo"`
	LastCommit *LastCommit `json:"last_commit"`
}

type ProtocolDetail struct {
	Slug        string              `json:"slug"`
	Name        string              `json:"name"`
	Category    *string             `json:"category,omitempty"`
	Chains      []string            `json:"chains"`
	Methodology map[string]string   `json:"methodology"`
	Dimensions  []ProtocolDimension `json:"dimensions"`
}
