package web

import "github.com/go-chi/chi/v5"

func DefaultRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/health", HealthcheckHandler)
	return r
}
