package apiserver

import (
	"context"
	"github.com/BellZaph/banner-roulette-backend/internal/apiserver/resources"
	v1 "github.com/BellZaph/banner-roulette-backend/internal/apiserver/v1"
	"github.com/BellZaph/banner-roulette-backend/internal/config"
	"github.com/BellZaph/banner-roulette-backend/internal/service"
	"github.com/BellZaph/banner-roulette-backend/pkg/hash"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)


const (
	filesDir      = "./assets"
	compressLevel = 5
)

type HTTPServer struct {
	Address           string
	FilesDir          string
	IsTesting         bool
	CertFile, KeyFile *string
	Services *service.Services
	Hasher hash.Hasher

	idleConnsClosed chan struct{}
	masterCtx       context.Context
	version         string
}

func NewHTTPServer(ctx context.Context, cfg *config.Config, services *service.Services, hasher hash.Hasher, version string) *HTTPServer {
	srv := &HTTPServer{
		Address:   ":" + cfg.HTTPPort,
		FilesDir:  filesDir,
		IsTesting: cfg.Debug,
		Services: services,
		Hasher: hasher,

		idleConnsClosed: make(chan struct{}),
		masterCtx:       ctx,
		version:         version,
	}

	if cfg.CertFile != "" {
		srv.CertFile = &cfg.CertFile
	}

	if cfg.KeyFile != "" {
		srv.KeyFile = &cfg.KeyFile
	}

	return srv
}


func allowedOrigins(testing bool) []string {
	if testing {
		return []string{"*"}
	}

	return []string{}
}

func (srv *HTTPServer) setupRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.NoCache)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.NewCompressor(compressLevel).Handler)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins(srv.IsTesting),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Mount("/v1", v1.Handler{Services: srv.Services, Hasher: srv.Hasher}.Routes())
	r.Mount("/assets", resources.FilesResource{FilesDir: filesDir}.Routes())

	return r
}

func (srv *HTTPServer) Run() error {
	const (
		readTimeout  = 5 * time.Second
		writeTimeout = 30 * time.Second
	)

	s := &http.Server{
		Addr:         srv.Address,
		Handler:      srv.setupRouter(),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	go srv.ListenCtxForGracefulTermination(s)
	logrus.Infof("serving HTTP on \"http://localhost%s\"", srv.Address)

	if srv.CertFile == nil && srv.KeyFile == nil {
		if err := s.ListenAndServe(); err != nil {
			return err
		}
	} else {
		if err := s.ListenAndServeTLS(*srv.CertFile, *srv.KeyFile); err != nil {
			return err
		}
	}

	return nil
}

func (srv *HTTPServer) ListenCtxForGracefulTermination(s *http.Server) {
	<-srv.masterCtx.Done()

	if err := s.Shutdown(srv.masterCtx); err != nil {
		logrus.Infof("HTTP server Shutdown: %v", err)
	}

	logrus.Println("Processed idle connections successfully before termination")
	close(srv.idleConnsClosed)
}

func (srv *HTTPServer) WaitForGracefulTermination() {
	<-srv.idleConnsClosed
}

