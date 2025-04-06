package company

import (
	"gobizmanager/pkg/language"
	"gobizmanager/platform/middleware/ratelimit"

	"net/http"

	"github.com/go-chi/chi/v5"
)

func Routes(handler *Handler, msgStore *language.MessageStore) http.Handler {
	r := chi.NewRouter()

	// Apply rate limiting middleware to all company routes
	r.Group(func(r chi.Router) {
		r.Use(ratelimit.New(100))
		r.Post("/", handler.createCompany)
		r.Get("/", handler.listCompanies)
		r.Get("/{id}", handler.getCompany)
		r.Put("/{id}", handler.updateCompany)
		r.Delete("/{id}", handler.deleteCompany)
	})

	return r
}
