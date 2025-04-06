package context

import (
	"context"
	"net/http"

	"gobizmanager/pkg/language"
)

type contextKey string

const (
	LanguageKey contextKey = "language"
)

// GetLanguage returns the language from the context
func GetLanguage(ctx context.Context) string {
	if lang, ok := ctx.Value(LanguageKey).(string); ok {
		return lang
	}
	return language.English // Default to English
}

// SetLanguage sets the language in the context
func SetLanguage(ctx context.Context, lang string) context.Context {
	return context.WithValue(ctx, LanguageKey, lang)
}

// LanguageMiddleware creates a middleware that sets the language based on the Accept-Language header
func LanguageMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get language from Accept-Language header
			lang := r.Header.Get("Accept-Language")
			if lang == "" {
				lang = language.English // Default to English
			}

			// Set language in context
			ctx := SetLanguage(r.Context(), lang)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
