package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"com.github.denisbytes.goimageai/db"
	"com.github.denisbytes.goimageai/types"
	"com.github.denisbytes.goimageai/view/generate"
	"github.com/go-chi/chi/v5"
)

func HandleGenerateIndex(w http.ResponseWriter, r *http.Request) error {
	user := GetAuthenticatedUser(r)
	images, err := db.GetImagesByUserID(user.ID)
	if err != nil {
		return err
	}
	data := generate.ViewData{
		Images: images,
	}
	return generate.Index(data).Render(r.Context(), w)
}

func HandleGeneratePost(w http.ResponseWriter, r *http.Request) error {
	user := GetAuthenticatedUser(r)
	prompt := "red sportscar in a garden"
	img := types.Image{
		Prompt: prompt,
		UserID: user.ID,
		Status: types.ImageStatusPending,
	}
	if err := db.CreateImage(&img); err != nil {
		return err
	}
	return generate.GalleryImage(img).Render(r.Context(), w)
}

func HandleGenerateImageStatus(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return err
	}
	image, err := db.GetImagesByID(id)
	if err != nil {
		return err
	}
	slog.Info("Checking images status", "id", id)
	return generate.GalleryImage(image).Render(r.Context(), w)
}
