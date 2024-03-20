package handler

import (
	"log/slog"
	"net/http"
	"os"

	"com.github.denisbytes.goimageai/db"
	"com.github.denisbytes.goimageai/pkg/kit/validate"
	"com.github.denisbytes.goimageai/pkg/sb"
	"com.github.denisbytes.goimageai/types"
	"com.github.denisbytes.goimageai/view/auth"
	"github.com/gorilla/sessions"
	"github.com/nedpals/supabase-go"
)

const (
	sessionUserKey        = "user"
	sessionAccessTokenKey = "accessToken"
)

func HandleLogInIndex(w http.ResponseWriter, r *http.Request) error {
	return auth.Login().Render(r.Context(), w)
}

func HandleLoginPost(w http.ResponseWriter, r *http.Request) error {
	credentials := supabase.UserCredentials{
		Email: r.FormValue("email"),
	}
	err := sb.Client.Auth.SendMagicLink(r.Context(), credentials.Email)
	if err != nil {
		slog.Error("login error", "err", err)
		return auth.LoginForm(credentials, auth.LoginErrors{
			InvalidCredentials: err.Error(),
		}).Render(r.Context(), w)
	}
	// if err := setAuthSession(w, r, resp.AccessToken); err != nil {
	// 	return err
	// }
	return auth.MagicLinkSuccess(credentials.Email).Render(r.Context(), w)
}

func HandleAuthCallback(w http.ResponseWriter, r *http.Request) error {
	accessToken := r.URL.Query().Get("access_token")
	if len(accessToken) == 0 {
		return auth.CallbackScript().Render(r.Context(), w)
	}
	if err := setAuthSession(w, r, accessToken); err != nil {
		return err
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func setAuthSession(w http.ResponseWriter, r *http.Request, accessToken string) error {
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	session, _ := store.Get(r, "user")
	session.Values["accessToken"] = accessToken
	return session.Save(r, w)
}

func HandleLogoutPost(w http.ResponseWriter, r *http.Request) error {
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	session, _ := store.Get(r, "user")
	session.Values["accessToken"] = ""
	session.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
	return nil
}

func HandleLoginWithGithubPost(w http.ResponseWriter, r *http.Request) error {
	resp, err := sb.Client.Auth.SignInWithProvider(supabase.ProviderSignInOptions{
		Provider:   "github",
		RedirectTo: "http://localhost:3000/auth/callback",
	})
	if err != nil {
		return err
	}
	http.Redirect(w, r, resp.URL, http.StatusSeeOther)
	return nil
}

func HandleAccountSetupIndex(w http.ResponseWriter, r *http.Request) error {
	return auth.AccountSetup().Render(r.Context(), w)
}

func HandleAccountSetupPost(w http.ResponseWriter, r *http.Request) error {
	params := auth.AccountSetupParams{
		Username: r.FormValue("username"),
	}
	var errors auth.AccountSetupErrors
	ok := validate.New(&params, validate.Fields{
		"Username": validate.Rules(validate.Min(2), validate.Max(50)),
	}).Validate(&errors)
	if !ok {
		return auth.AccountSetupForm(params, errors).Render(r.Context(), w)
	}
	user := GetAuthenticatedUser(r)
	account := types.Account{

		UserID:   user.ID,
		Username: params.Username,
	}
	if err := db.CreateAccount(&account); err != nil {
		return err
	}
	return hxRedirect(w, r, "/")
}
