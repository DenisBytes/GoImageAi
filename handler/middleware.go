package handler

import (
	"context"
	"net/http"
	"os"
	"strings"

	"com.github.denisbytes.goimageai/pkg/sb"
	"com.github.denisbytes.goimageai/types"
	"github.com/gorilla/sessions"
)

func WithUser(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}

		store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
		session, err := store.Get(r, "user")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		accessToken := session.Values["accessToken"]
		resp, err := sb.Client.Auth.User(r.Context(), accessToken.(string))
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		user := types.AuthenticatedUser{
			Email:    resp.Email,
			LoggedIn: true,
		}

		ctx := context.WithValue(r.Context(), types.UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func WithAuth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}
		user := GetAuthenticatedUser(r)
		if !user.LoggedIn {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
