package v1

import (
	"github.com/BellZaph/banner-roulette-backend/internal/service"
	"github.com/BellZaph/banner-roulette-backend/pkg/hash"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Services *service.Services
	Hasher hash.Hasher
}

func (h Handler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Mount("/banners", h.initBannersRoutes())

	return r
}
