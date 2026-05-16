// GET /api/matrix handler. Returns the pinned columns plus filtered protocol rows.
package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
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
	_, err := ParseMatrixQuery(r)
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
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(map[string]any{"error": map[string]string{"code": code, "message": message}})
		return
	}
	// Cells and rows population lands in a follow-up step; for now return the
	// pinned columns with an empty rows slice once the query has validated.
	resp := MatrixResponse{
		Columns: ColumnList(),
		Rows:    []Row{},
		Total:   0,
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("matrix encode failed", "error", err)
	}
}
