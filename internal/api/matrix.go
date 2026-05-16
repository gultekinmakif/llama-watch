// GET /api/matrix handler. Returns the pinned columns plus filtered protocol rows.
package api

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gultekinmakif/llama-watch/internal/db/postgres"
	"github.com/gultekinmakif/llama-watch/internal/registry"
)

// Matrix serves GET /api/matrix.
func Matrix(w http.ResponseWriter, r *http.Request) {
	q, err := ParseMatrixQuery(r)
	if err != nil {
		status := http.StatusInternalServerError
		code := "internal"
		message := "internal error"
		var perr *ParseError
		if errors.As(err, &perr) {
			status = http.StatusBadRequest
			code = perr.Code
			message = perr.Message
		}
		writeErr(w, status, code, message)
		return
	}

	db := postgres.Get()
	ctx := r.Context()

	total, err := countProtocols(ctx, db)
	if err != nil {
		slog.Error("matrix count failed", "error", err)
		writeErr(w, http.StatusInternalServerError, "internal", "internal error")
		return
	}

	rows, err := listProtocols(ctx, db, q.Limit, q.Offset)
	if err != nil {
		slog.Error("matrix list failed", "error", err)
		writeErr(w, http.StatusInternalServerError, "internal", "internal error")
		return
	}

	writeJSON(w, http.StatusOK, MatrixResponse{
		Columns: registry.Columns(),
		Rows:    rows,
		Total:   total,
	})
}

type Row struct {
	Slug     string         `json:"slug"`
	Name     string         `json:"name"`
	Category *string        `json:"category,omitempty"`
	TVL      *float64       `json:"tvl,omitempty"`
	Chains   []string       `json:"chains"`
	Cells    map[string]int `json:"cells"`
}

type MatrixResponse struct {
	Columns []registry.Column `json:"columns"`
	Rows    []Row             `json:"rows"`
	Total   int               `json:"total"`
}
