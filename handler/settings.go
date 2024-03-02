package handler

import (
	"net/http"

	"com.github.denisbytes.goimageai/view/settings"
)

func HandleSettingsIndex(w http.ResponseWriter, r *http.Request) error {
	user := GetAuthenticatedUser(r)
	return settings.Index(user).Render(r.Context(), w)
}
