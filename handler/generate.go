package handler

import (
	"log/slog"
	"net/http"

	"com.github.denisbytes.goimageai/types"
	"com.github.denisbytes.goimageai/view/generate"
	"github.com/go-chi/chi/v5"
)

func HandleGenerateIndex(w http.ResponseWriter, r *http.Request) error {
	images := make([]types.Image, 20)
	data := generate.ViewData{
		Images: images,
	}
	images[0].Status = types.ImageStatusPending
	return generate.Index(data).Render(r.Context(), w)
}

func HandleGeneratePost(w http.ResponseWriter, r *http.Request) error {
	return generate.GalleryImage(types.Image{Status: types.ImageStatusPending}).Render(r.Context(), w)
}

func HandleGenerateImageStatus(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	image := types.Image{
		Status: types.ImageStatusPending,
	}	
	slog.Info("Checking images status", "id", id)
	return generate.GalleryImage(image).Render(r.Context(), w)
}
