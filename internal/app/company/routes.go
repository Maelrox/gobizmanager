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
		r.Post("/", handler.CreateCompany)
		r.Get("/", handler.ListCompanies)
		r.Get("/{id}", handler.GetCompany)
		r.Put("/{id}", handler.UpdateCompany)
		r.Delete("/{id}", handler.DeleteCompany)
	})

	return r
}
