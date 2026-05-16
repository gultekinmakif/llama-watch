// GET /health handler. Pings the DB with a 2s timeout and reports latency.
package api

import (
	"context"
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

// Health returns 200 with status:ok when the DB ping succeeds within 2 seconds,
// otherwise 503 with status:down. db_ms reports ping latency either way.
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
