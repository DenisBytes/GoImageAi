package handler

import (
	"net/http"

	"com.github.denisbytes.goimageai/view/auth"
)

func HandleLogInIndex(w http.ResponseWriter, r *http.Request) error {
	return auth.LogIn().Render(r.Context(), w)
}
