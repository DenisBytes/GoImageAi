package handler

import (
	"net/http"

	"com.github.denisbytes.goimageai/db"
	"com.github.denisbytes.goimageai/pkg/kit/validate"
	"com.github.denisbytes.goimageai/view/settings"
)

func HandleSettingsIndex(w http.ResponseWriter, r *http.Request) error {
	user := GetAuthenticatedUser(r)
	return settings.Index(user).Render(r.Context(), w)
}

func HandleSettingsUsernameUpdate(w http.ResponseWriter, r *http.Request) error {
	params := settings.ProfileParams{
		Username: r.FormValue("username"),
	}
	errors := settings.ProfileErrors{}
	ok := validate.New(&params, validate.Fields{
		"Username": validate.Rules(validate.Min(3), validate.Max(40)),
	}).Validate(&errors)
	if !ok {
		return settings.ProfileForm(params, errors).Render(r.Context(), w)
	}
	user := GetAuthenticatedUser(r)
	user.Account.Username = params.Username
	if err := db.UpdateAccount(&user.Account); err != nil {
		return err
	}
	params.Success = true
	return settings.ProfileForm(params, settings.ProfileErrors{}).Render(r.Context(), w)
}
