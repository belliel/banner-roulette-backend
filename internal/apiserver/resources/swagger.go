package resources

import (
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"path/filepath"
)

type SwaggerResource struct {
	FilesPath string
}

func (sr SwaggerResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/*", httpSwagger.Handler(
		httpSwagger.URL(filepath.Join(sr.FilesPath, "swagger.json")),
	))

	filesRoot := http.Dir(sr.FilesPath)

	NewFileServer(r, "/docs", filesRoot)

	return r
}