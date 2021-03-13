package v1

import (
	"encoding/json"
	"fmt"
	"github.com/BellZaph/banner-roulette-backend/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

func (h Handler) initBannersRoutes() chi.Router {
	r := chi.NewRouter()

	r.Use(ForceJSONContentType())
	r.Get("/banner-{bannerId}", h.getBannerById)
	r.Get("/random", h.getBannerRandom)
	r.Get("/randoms", h.getBannerRandoms)
	r.Get("/", h.getByPage)
	r.Put("/banner-{bannerId}", h.incrementByID)
	r.Put("/", h.putBannerById)
	r.Post("/", h.createBanner)
	r.Post("/images/upload", h.uploadImage)
	r.Delete("/{bannerId}", h.deleteBanner)


	return r
}

func (h Handler) getBannerById(w http.ResponseWriter, r *http.Request) {
	bannerIdString := chi.URLParam(r, "bannerId")

	bannerId, err := uuid.Parse(bannerIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}


	banner, err := h.Services.Banners.GetById(r.Context(), bannerId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	_ = json.NewEncoder(w).Encode(banner)
}

func (h Handler) putBannerById(w http.ResponseWriter, r *http.Request) {
	var banner service.BannerUpdateInput

	if err := json.NewDecoder(r.Body).Decode(&banner); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.Services.Banners.Update(r.Context(), banner); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}
}

func (h Handler) createBanner(w http.ResponseWriter, r *http.Request) {
	var banner service.BannerCreateInput

	if err := json.NewDecoder(r.Body).Decode(&banner); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.Services.Banners.Create(r.Context(), banner); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(SuccessResponse{Success: true})
}

func (h Handler) deleteBanner(w http.ResponseWriter, r *http.Request) {
	bannerIdString := chi.URLParam(r, "bannerId")

	bannerId, err := uuid.Parse(bannerIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	if err = h.Services.Banners.Delete(r.Context(), bannerId); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	_ = json.NewEncoder(w).Encode(SuccessResponse{Success: true})
}

func (h Handler) getBannerRandom(w http.ResponseWriter, r *http.Request) {
	hourString := r.URL.Query().Get("hour")

	hour, err := strconv.Atoi(hourString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}


	banner, err := h.Services.Banners.GetRandom(r.Context(), hour)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	_ = json.NewEncoder(w).Encode(banner)
}

func (h Handler) getBannerRandoms(w http.ResponseWriter, r *http.Request) {
	hourString := r.URL.Query().Get("hour")

	hour, err := strconv.Atoi(hourString)
	if err != nil {
		return
	}


	banner, err := h.Services.Banners.GetRandoms(r.Context(), hour, 2)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	_ = json.NewEncoder(w).Encode(banner)
}

func (h Handler) getByPage(w http.ResponseWriter, r *http.Request) {
	pageString := r.URL.Query().Get("page")

	page, err := strconv.Atoi(pageString)
	if err != nil {
		page = 0
	}

	banners, err := h.Services.Banners.GetByPage(r.Context(), page)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	_ = json.NewEncoder(w).Encode(banners)
}

func (h Handler) incrementByID(w http.ResponseWriter, r *http.Request) {
	bannerIdString := chi.URLParam(r, "bannerId")

	bannerId, err := uuid.Parse(bannerIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}


	if err := h.Services.Banners.IncrementCount(r.Context(), bannerId); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	_ = json.NewEncoder(w).Encode(SuccessResponse{Success: true})
}

func (h Handler) uploadImage(w http.ResponseWriter, r *http.Request) {

	file, handler, err := r.FormFile("image")
	if err != nil {
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	extension := strings.Split(handler.Filename, ".")[len(strings.Split(handler.Filename, "."))-1]

	fileName := h.Hasher.Hash()

	filePath := path.Join("assets/images", fmt.Sprintf("%s.%s", fileName, extension))

	dst, err := os.Create(filePath)
	if err != nil {
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(ImageUploadResponse{ImageURI: filePath})
}
