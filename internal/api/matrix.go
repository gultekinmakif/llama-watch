// GET /api/matrix handler. Returns the pinned columns plus filtered protocol rows.
package api

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gultekinmakif/llama-watch/internal/db/postgres"
)

type Column struct {
	Key   string `json:"key"`
	Label string `json:"label"`
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
	Columns []Column `json:"columns"`
	Rows    []Row    `json:"rows"`
	Total   int      `json:"total"`
}

// Closed, fixed column set. Order is load-bearing. Unexported so callers cannot
// mutate the backing array; use ColumnList for a safe copy.
var columns = []Column{
	{Key: "tvl", Label: "TVL"},
	{Key: "dailyFees", Label: "Daily Fees"},
	{Key: "dailyRevenue", Label: "Daily Revenue"},
	{Key: "dailyVolume", Label: "Daily Volume"},
	{Key: "dailyNotionalVolume", Label: "Notional Volume"},
	{Key: "dailyPremiumVolume", Label: "Premium Volume"},
	{Key: "openInterestAtEnd", Label: "Open Interest"},
	{Key: "dailyBridgeVolume", Label: "Bridge Volume"},
	{Key: "dailyActiveUsers", Label: "Active Users"},
}

// ColumnList returns a copy of the pinned column set so external callers can read
// the matrix shape without risking mutation of the package-level slice.
func ColumnList() []Column {
	out := make([]Column, len(columns))
	copy(out, columns)
	return out
}

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
		writeJSON(w, status, map[string]any{"error": map[string]string{"code": code, "message": message}})
		return
	}

	db := postgres.Get()
	ctx := r.Context()

	total, err := countProtocols(ctx, db)
	if err != nil {
		slog.Error("matrix count failed", "error", err)
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": map[string]string{"code": "internal", "message": "internal error"}})
		return
	}

	rows, err := listProtocols(ctx, db, q.Limit, q.Offset)
	if err != nil {
		slog.Error("matrix list failed", "error", err)
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": map[string]string{"code": "internal", "message": "internal error"}})
		return
	}

	writeJSON(w, http.StatusOK, MatrixResponse{
		Columns: ColumnList(),
		Rows:    rows,
		Total:   total,
	})
}
