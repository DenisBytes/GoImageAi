package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"com.github.denisbytes.goimageai/db"
	"com.github.denisbytes.goimageai/types"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ReplicateResponse struct {
	Status string   `json:"status"`
	Output []string `json:"output"`
}

func HandleReplicateCallback(w http.ResponseWriter, r *http.Request) error {
	var resp ReplicateResponse

	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return err
	}
	if resp.Status != "succeeded" {
		return fmt.Errorf("replicate callback responded with a non ok status: %s", resp.Status)
	}

	batchID, err := uuid.Parse(chi.URLParam(r, "batchID"))
	if err != nil {
		return fmt.Errorf("replicate callback invalud batchID %s", err)
	}

	images, err := db.GetImagesByBatchID(batchID)
	if err != nil {
		return fmt.Errorf("replicate callback failed to find image with batchID %s. err: %s", batchID, err)
	}

	if len(images) != len(resp.Output) {
		return fmt.Errorf("replicate callback unequal images compared to replicate outputs")
	}

	for i, imageURL := range resp.Output {
		images[i].Status = types.ImageStatusCompleted
		images[i].ImageLocation = imageURL
	}

	return nil
}
