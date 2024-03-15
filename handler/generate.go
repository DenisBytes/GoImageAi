package handler

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"com.github.denisbytes.goimageai/db"
	"com.github.denisbytes.goimageai/pkg/kit/validate"
	"com.github.denisbytes.goimageai/types"
	"com.github.denisbytes.goimageai/view/generate"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/replicate/replicate-go"
	"github.com/uptrace/bun"
)

const creditsPerImage = 2

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

	if amount <= 0 || amount > 8 {
		errors.Amount = "Please enter a valid amount"
		return generate.Form(params, errors).Render(r.Context(), w)
	}

	ok := validate.New(params, validate.Fields{
		"Prompt": validate.Rules(validate.Min(10), validate.Max(100)),
	}).Validate(&errors)
	if !ok {
		return generate.Form(params, errors).Render(r.Context(), w)
	}

	creditsNeeded := params.Amount * creditsPerImage
	if user.Account.Credits < creditsNeeded {
		errors.CreditsNeeded = creditsNeeded
		errors.UserCredits = user.Account.Credits
		errors.Credits = true
		return generate.Form(params, errors).Render(r.Context(), w)
	}

	batchID := uuid.New()
	genParams := GenerateImageParams{
		Prompt:  params.Prompt,
		Amount:  params.Amount,
		UserID:  user.ID,
		BatchID: batchID,
	}

	if err := generateImages(r.Context(), genParams); err != nil {
		return err
	}

	err := db.Bun.RunInTx(r.Context(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		for i := 0; i < params.Amount; i++ {
			img := types.Image{
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
}

func HandleGenerateImageStatus(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return err
	}
	image, err := db.GetImageByID(id)
	if err != nil {
		return err
	}
	slog.Info("checking image status", "id", id)
	return generate.GalleryImage(image).Render(r.Context(), w)
}

type GenerateImageParams struct {
	Prompt  string
	Amount  int
	BatchID uuid.UUID
	UserID  uuid.UUID
}

func generateImages(ctx context.Context, params GenerateImageParams) error {
	r8, err := replicate.NewClient(replicate.WithTokenFromEnv())
	if err != nil {
		log.Fatal(err)
	}
	input := replicate.PredictionInput{
		"prompt":      params.Prompt,
		"num_outputs": params.Amount,
	}
	webhook := replicate.Webhook{
		URL:    fmt.Sprintf("https://webhook.site/5fe9f245-ef5c-432c-931e-6eeff587066b/%s/%s", params.UserID, params.BatchID),
		Events: []replicate.WebhookEventType{"start", "completed"},
	}
	_, err = r8.CreatePrediction(ctx, "d70beb400d223e6432425a5299910329c6050c6abcf97b8c70537d6a1fcb269a", input, &webhook, false)
	return err
}
