package handler

import (
	"net/http"

	"com.github.denisbytes.goimageai/view/home"
)

func HandleHomeIndex(w http.ResponseWriter, r *http.Request) error {
	return home.Index().Render(r.Context(), w)
}
