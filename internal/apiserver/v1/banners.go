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

// getBannerById godoc
// @Summary Banner by ID
// @Description Get banner by UUID id in http param
// @ID get-banner-by-id
// @Accept html
// @Produce json
// @Param bannerId path string true "bannerId"
// @Success 200 {object} models.Banner
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /banners/banner-{bannerId} [get]
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

// putBannerById godoc
// @Summary Update Banner by ID
// @Description Update banner by UUID id in http param
// @ID put-banner-by-id
// @Accept json
// @Produce json
// @Param banner body service.BannerUpdateInput true "bannerInput"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /banners [put]
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

	_ = json.NewEncoder(w).Encode(SuccessResponse{Success: true})
}

// createBanner godoc
// @Summary Create banner
// @Description Create new banner
// @ID create-banner
// @Accept json
// @Produce json
// @Param banner body service.BannerCreateInput true "bannerInput"
// @Success 201 {object} SuccessResponse
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /banners [post]
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


// deleteBanner godoc
// @Summary Delete banner
// @Description Delete banner by id
// @ID delete-banner
// @Accept html
// @Produce json
// @Param bannerId path string true "bannerId"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /banners/{bannerId} [delete]
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

// getBannerRandom godoc
// @Summary Get random banner
// @Description Get random banner from database
// @ID get-banner-random
// @Accept html
// @Produce json
// @Param hour query int true "hour"
// @Success 200 {object} models.Banner
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /banners/random [get]
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

// getBannerRandoms godoc
// @Summary Get random banners with limit
// @Description Get random banners from database
// @ID get-banner-randoms
// @Accept html
// @Produce json
// @Param hour query int true "hour"
// @Success 200 {array} models.Banner
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /banners/randoms [get]
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

// getByPage godoc
// @Summary Get ordered list of banners
// @Description Get ordered list of banners
// @ID get-by-page
// @Accept html
// @Produce json
// @Param page query int true "page"
// @Success 200 {array} models.Banner
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /banners [get]
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

// incrementByID godoc
// @Summary Adds 1 show counter
// @Description And its all
// @ID increment-by-id
// @Accept html
// @Produce json
// @Param bannerId path int true "bannerId"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /banner-{bannerId} [put]
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

// uploadImage godoc
// @Summary Hand file from browser
// @Description multipart form data, i dont know how to pass it to Params
// @ID upload-image
// @Accept multipart/form-data
// @Produce json
// @Success 200 {object} ImageUploadResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /banners/images/upload [post]
func (h Handler) uploadImage(w http.ResponseWriter, r *http.Request) {

	file, handler, err := r.FormFile("image")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}
	defer file.Close()

	extension := strings.Split(handler.Filename, ".")[len(strings.Split(handler.Filename, "."))-1]

	fileName := h.Hasher.Hash()

	filePath := path.Join("assets/images", fmt.Sprintf("%s.%s", fileName, extension))

	dst, err := os.Create(filePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	_ = json.NewEncoder(w).Encode(ImageUploadResponse{ImageURI: filePath})
}
