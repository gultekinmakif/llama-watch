// GET /api/matrix handler. Returns the pinned columns plus filtered protocol rows.
package api

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gultekinmakif/llama-watch/internal/db/postgres"
	"github.com/gultekinmakif/llama-watch/internal/registry"
	"gorm.io/gorm"
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

	var total int
	var rows []Row
	if err := db.Transaction(func(tx *gorm.DB) error {
		var ierr error
		total, ierr = countProtocols(ctx, tx, q)
		if ierr != nil {
			slog.Error("matrix count failed", "error", ierr)
			return ierr
		}
		rows, ierr = listProtocols(ctx, tx, q)
		if ierr != nil {
			slog.Error("matrix list failed", "error", ierr)
		}
		return ierr
	}); err != nil {
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
	Slug     string                        `json:"slug"`
	Name     string                        `json:"name"`
	Category *string                       `json:"category,omitempty"`
	TvlUSD   *float64                      `json:"tvl_usd,omitempty"`
	Chains   []string                      `json:"chains"`
	Cells    map[string]registry.CellState `json:"cells"`
}

type MatrixResponse struct {
	Columns []registry.Column `json:"columns"`
	Rows    []Row             `json:"rows"`
	Total   int               `json:"total"`
}
