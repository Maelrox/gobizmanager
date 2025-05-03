package auth

import (
	"net/http"

	"gobizmanager/pkg/context"
	"gobizmanager/pkg/language"
	"gobizmanager/pkg/sanitizer"
	"gobizmanager/platform/middleware/ratelimit"

	"github.com/go-chi/chi/v5"
)

func Routes(handler *Handler, msgStore *language.MessageStore) http.Handler {
	r := chi.NewRouter()
	r.Use(context.LanguageMiddleware())
	r.Use(sanitizer.SanitizeInput)
	r.Use(sanitizer.SanitizeOutput)
	r.With(ratelimit.New(10)).Post("/login", handler.Login)
	r.Post("/register", handler.Register)
	r.Post("/refresh", handler.RefreshToken)
	return r
}
