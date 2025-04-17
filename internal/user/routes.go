package user

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Routes(handler *Handler) http.Handler {
	r := chi.NewRouter()

	r.Get("/search", handler.SearchUsers)

	return r
}
