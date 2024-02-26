package main

import (
	"embed"
	"log"
	"log/slog"
	"net/http"
	"os"

	"com.github.denisbytes.goimageai/handler"
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
	port := os.Getenv("HTTP_LISTEN_ADDR")
	slog.Info("application running", "port", port)
	log.Fatal(http.ListenAndServe(port, router))
}

func initEnvVar() error {
	return godotenv.Load()
}
