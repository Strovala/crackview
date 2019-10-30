package http

import (
	"net/http"

	"github.com/go-chi/chi"
	chiCors "github.com/go-chi/cors"
)

// NewCrackviewHandler creates new application handler
func NewCrackviewHandler() http.Handler {
	cors := chiCors.New(chiCors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	})

	api := chi.NewMux()
	api.Use(cors.Handler)
	api.Use(JSONRecoverer)

	api.Mount("/", newCodeHandler())
	return api
}
