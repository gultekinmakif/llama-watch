// GET /api/chains handler. Returns the distinct set of chains across protocols with per-chain protocol counts.
package api

import (
	"log/slog"
	"net/http"

	"github.com/gultekinmakif/llama-watch/internal/db/postgres"
)

type ChainEntry struct {
	Key           string `json:"key"`
	Label         string `json:"label"`
	ProtocolCount int    `json:"protocol_count"`
}

type chainsResponse struct {
	Chains []ChainEntry `json:"chains"`
}

// Chains serves GET /api/chains.
func Chains(w http.ResponseWriter, r *http.Request) {
	db := postgres.Get()
	ctx := r.Context()

	rows, err := listChains(ctx, db)
	if err != nil {
		slog.Error("chains list failed", "error", err)
		writeErr(w, http.StatusInternalServerError, "internal", "internal error")
		return
	}

	writeJSON(w, http.StatusOK, chainsResponse{Chains: rows})
}
