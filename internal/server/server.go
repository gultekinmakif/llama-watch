// 2.2: HTTP server bootstrap. Mounts routes, applies middleware, handles graceful shutdown.
// Middleware chain at request time: 2.3 recoverer then 2.4 request_id then 2.5 logger then 2.6 handler.
package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/gultekinmakif/llama-watch/internal/api"
	"github.com/gultekinmakif/llama-watch/internal/config"
	"github.com/gultekinmakif/llama-watch/internal/middleware"
)

type Server struct {
	config *config.Config
	server *http.Server
}

func New(cfg *config.Config) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", api.Health)
	mux.HandleFunc("GET /api/chains", api.Chains)
	mux.HandleFunc("GET /api/matrix", api.Matrix)
	mux.HandleFunc("GET /api/matrix/{slug}", api.MatrixDetail)
	// Static export from `web/out/`. Returns 404 from the FS when the path is
	// missing; no SPA fallback to index.html so unknown routes 404 cleanly.
	mux.Handle("/", http.FileServer(http.Dir("web/out")))

	return &Server{
		config: cfg,
		server: &http.Server{
			Addr:         ":" + cfg.Port,
			Handler:      middleware.Recoverer(middleware.RequestID(middleware.Logger(mux))),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}

func (s *Server) Start(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		slog.Info("server listening", "port", s.config.Port, "env", s.config.Env)
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		slog.Info("shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
		defer cancel()
		return s.server.Shutdown(shutdownCtx)
	case err := <-errCh:
		return err
	}
}
