package handler

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"com.github.denisbytes.goimageai/db"
	"com.github.denisbytes.goimageai/pkg/kit/validate"
	"com.github.denisbytes.goimageai/types"
	"com.github.denisbytes.goimageai/view/generate"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
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

	amount, _ := strconv.Atoi(r.FormValue("amount"))
	params := generate.FormData{
		Prompt: r.FormValue("prompt"),
		Amount: amount,
	}
	var errors generate.FormErrors

	fmt.Println(params.Prompt)

	if amount <= 0 || amount > 8 {
		errors.Amount = "Please enter a valid amount"
		return generate.Form(params, errors).Render(r.Context(), w)
	}
	ok := validate.New(params, validate.Fields{
		"Prompt": validate.Rules(validate.Max(100), validate.Min(10)),
	}).Validate(&errors)

	if !ok {
		return generate.Form(params, errors).Render(r.Context(), w)
	}

	err := db.Bun.RunInTx(r.Context(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		batchID := uuid.New()
		for i := 0; i < params.Amount; i++ {
			img := types.Image{
				Prompt:  params.Prompt,
				UserID:  user.ID,
				Status:  types.ImageStatusPending,
				BatchID: batchID,
			}
			if err := db.CreateImage(&img); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return hxRedirect(w, r, "/generate")
	// return generate.GalleryImage(img).Render(r.Context(), w)
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
