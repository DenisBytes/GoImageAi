package handler

import (
	"log/slog"
	"net/http"

	"com.github.denisbytes.goimageai/pkg/sb"
	"com.github.denisbytes.goimageai/pkg/util"
	"com.github.denisbytes.goimageai/view/auth"
	"github.com/nedpals/supabase-go"
)

func HandleLogInIndex(w http.ResponseWriter, r *http.Request) error {
	return auth.Login().Render(r.Context(), w)
}

func HandleLoginPost(w http.ResponseWriter, r *http.Request) error {
	credentials := supabase.UserCredentials{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	if !util.IsValidEmail(credentials.Email) {
		return auth.LoginForm(credentials, auth.LoginErrors{
			Email:    "Please Enter a valid email",
			Password: "",
		}).Render(r.Context(), w)
	}

	resp, err := sb.Client.Auth.SignIn(r.Context(), credentials)
	if err != nil {
		slog.Error("login error", "err", err)
		return auth.LoginForm(credentials, auth.LoginErrors{
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}).Render(r.Context(), w)
	}

	cookie := &http.Cookie{
		Value:    resp.AccessToken,
		Name:     "at",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)

	return nil
}
