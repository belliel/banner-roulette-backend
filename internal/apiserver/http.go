package apiserver

import (
	"context"
	"github.com/BellZaph/banner-roulette-backend/internal/config"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config, handler http.Handler) *Server {

	return &Server{
		httpServer: &http.Server{
			Addr:           ":" + cfg.HTTPPort,
			Handler:        handler,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
