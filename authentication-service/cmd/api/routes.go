package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	//specify who is allowed to connect

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"POST", "GET", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Acccept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Links"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	return mux
}
