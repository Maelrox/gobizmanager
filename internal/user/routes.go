package user

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Routes(handler *Handler) http.Handler {
	r := chi.NewRouter()

	// Add user search route
	r.Get("/search", handler.SearchUsers)

	return r
}
