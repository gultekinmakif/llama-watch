// HTTP handlers outside the /api surface. Only the health probe lives here.
package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/gultekinmakif/llama-watch/internal/db/postgres"
)

type healthResponse struct {
	Status string `json:"status"`
	DB     string `json:"db"`
	DBMs   int64  `json:"db_ms"`
}

func writeHeaders(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	writeHeaders(w, status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("json encode failed", "error", err)
	}
}

// Health returns 200 with status:ok when the DB ping succeeds within 2 seconds.
// otherwise 503 with status:down.
// db_ms reports the ping latency either way (timeout duration on failure).
func Health(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	start := time.Now()
	err := postgres.Ping(ctx)
	elapsed := time.Since(start).Milliseconds()

	resp := healthResponse{Status: "ok", DB: "ok", DBMs: elapsed}
	code := http.StatusOK
	if err != nil {
		slog.Warn("health: db ping failed", "error", err, "elapsed_ms", elapsed)
		resp.Status = "down"
		resp.DB = "down"
		code = http.StatusServiceUnavailable
	}

	writeJSON(w, code, resp)
}
