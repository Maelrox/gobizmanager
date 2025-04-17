package company_user

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Routes(handler *Handler) http.Handler {
	r := chi.NewRouter()

	r.Post("/register", handler.RegisterCompanyUser)

	return r
}
