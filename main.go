package main

import (
	"embed"
	"log"
	"log/slog"
	"net/http"
	"os"

	"com.github.denisbytes.goimageai/db"
	"com.github.denisbytes.goimageai/handler"
	"com.github.denisbytes.goimageai/pkg/sb"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

//go:embed public/*
var FS embed.FS

func main() {
	if err := initEnvVar(); err != nil {
		log.Fatal("Init err: ", err)
	}
	router := chi.NewMux()

	router.Use(handler.WithUser)

	// Handler for static files
	router.Handle("/*", http.StripPrefix("/", http.FileServer(http.FS(FS))))

	router.Get("/", handler.MakeHandler(handler.HandleHomeIndex))
	router.Get("/login", handler.MakeHandler(handler.HandleLogInIndex))
	router.Get("/login/provider/github", handler.MakeHandler(handler.HandleLoginWithGithubPost))
	router.Post("/login", handler.MakeHandler(handler.HandleLoginPost))
	router.Get("/signup", handler.MakeHandler(handler.HandleSignUpIndex))
	router.Post("/signup", handler.MakeHandler(handler.HandleSignUpPost))
	router.Get("/auth/callback", handler.MakeHandler(handler.HandleAuthCallback))
	router.Post("/logout", handler.MakeHandler(handler.HandleLogoutPost))
	router.Get("/account/setup", handler.MakeHandler(handler.HandleAccountSetupIndex))
	router.Post("/account/setup", handler.MakeHandler(handler.HandleAccountSetupPost))

	router.Group(func(auth chi.Router) {
		auth.Use(handler.WithAccountSetup)
		auth.Get("/settings", handler.MakeHandler(handler.HandleSettingsIndex))
		auth.Put("/settings/account/profile", handler.MakeHandler(handler.HandleSettingsUsernameUpdate))

		auth.Get("/auth/reset-password", handler.MakeHandler(handler.HandleResetPasswordIndex))
		auth.Post("/auth/reset-password", handler.MakeHandler(handler.HandleResetPasswordPost))
		auth.Put("/auth/reset-password", handler.MakeHandler(handler.HandleResetPasswordUpdate))

	})

	port := os.Getenv("HTTP_LISTEN_ADDR")
	slog.Info("application running", "port", port)
	log.Fatal(http.ListenAndServe(port, router))
}

func initEnvVar() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	if err := db.Init(); err != nil {
		return err
	}
	return sb.Init()
}
