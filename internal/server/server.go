package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/gultekinmakif/go-http-server/internal/config"
	"github.com/gultekinmakif/go-http-server/internal/handlers"
	"github.com/gultekinmakif/go-http-server/internal/middleware"
)

type Server struct {
	config *config.Config
	server *http.Server
}

func New(cfg *config.Config) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", handlers.Root)
	mux.HandleFunc("GET /health", handlers.Health)

	mux.HandleFunc("POST /posts", handlers.CreatePost)
	mux.HandleFunc("GET /posts", handlers.GetAllPost)
	mux.HandleFunc("GET /posts/{slug}", handlers.GetPost)

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
