package http

import (
	"encoding/json"
	"github.com/BellZaph/banner-roulette-backend/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
	"time"
)

type Handler struct {
	bannersService service.Banners
	services          *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{
		services:     services,
	}
}

func (h *Handler) Init(host, port string) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))


	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"pong": "done",
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	})

	return r
}

func (h *Handler) initAPI() chi.Router {
	return nil
}
