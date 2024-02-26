package handler

import (
	"context"
	"net/http"
	"strings"

	"com.github.denisbytes.goimageai/models"
)

func WithUser(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}
		user := models.AuthenticatedUser{
			Email: "agg@gmail.com",
			LoggedIn: true,
		}
		ctx := context.WithValue(r.Context(), models.UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
