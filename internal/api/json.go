// JSON response helpers shared across /api handlers.
package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// writeJSONHeader sets the JSON content type and writes the status line.
// Use when the caller wants to own the body write.
func writeJSONHeader(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
}

// writeJSON sets the JSON content type, writes the status, and encodes body.
// Encode errors are logged; the response status is already on the wire by then.
func writeJSON(w http.ResponseWriter, status int, body any) {
	writeJSONHeader(w, status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		slog.Error("json encode failed", "status", status, "error", err)
	}
}
