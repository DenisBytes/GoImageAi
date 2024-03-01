package handler

import (
	"log/slog"
	"net/http"

	"com.github.denisbytes.goimageai/pkg/kit/validate"
	"com.github.denisbytes.goimageai/pkg/sb"
	"com.github.denisbytes.goimageai/view/auth"
	"github.com/nedpals/supabase-go"
)

func HandleLogInIndex(w http.ResponseWriter, r *http.Request) error {
	return auth.Login().Render(r.Context(), w)
}

func HandleSignUpIndex(w http.ResponseWriter, r *http.Request) error {
	return auth.SignUp().Render(r.Context(), w)
}

func HandleSignUpPost(w http.ResponseWriter, r *http.Request) error {
	params := auth.SignUpParams{
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		ConfirmPassword: r.FormValue("confirmPassword"),
	}
	errors := auth.SignUpErrors{}
	if ok := validate.New(&params, validate.Fields{
		"Email":           validate.Rules(validate.Email),
		"Password":        validate.Rules(validate.Password),
		"ConfirmPassword": validate.Rules(validate.Equal(params.Password), validate.Message("Password don't match")),
	}).Validate(&errors); !ok {
		return auth.SignUpForm(params, errors).Render(r.Context(), w)
	}
	user, err := sb.Client.Auth.SignUp(r.Context(), supabase.UserCredentials{
		Email:    params.Email,
		Password: params.Password,
	})
	if err != nil {
		return err
	}
	return auth.SignUpSuccess(user.Email).Render(r.Context(), w)
}

func HandleLoginPost(w http.ResponseWriter, r *http.Request) error {
	credentials := supabase.UserCredentials{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
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
