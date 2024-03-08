package handler

import (
	"net/http"

	"com.github.denisbytes.goimageai/view/generate"
)

func HandleGenerateIndex(w http.ResponseWriter, r *http.Request) error {
	return generate.Index().Render(r.Context(), w)
}
