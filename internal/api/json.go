// JSON response helpers shared across /api handlers.
package api

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
)

// Hard-coded fallback envelope used when encoding the real body fails.
// Kept as a byte literal so we never re-enter the encode path while reporting an encode failure.
var encodeFailureBody = []byte(`{"error":{"code":"internal","message":"encode failed"}}` + "\n")

// writeJSON encodes body into a buffer first so a mid-stream encode error
// cannot leave a 2xx status with a truncated body on the wire.
func writeJSON(w http.ResponseWriter, status int, body any) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		slog.Error("json encode failed", "status", status, "error", err)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(encodeFailureBody)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(buf.Bytes())
}

// writeErr writes the {error:{code,message}} envelope shared by /api handlers.
func writeErr(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]any{"error": map[string]string{"code": code, "message": message}})
}
