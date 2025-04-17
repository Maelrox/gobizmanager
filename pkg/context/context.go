package context

import (
	"context"
	"net/http"

	"gobizmanager/pkg/language"
)

type contextKey string

const (
	LanguageKey     contextKey = "language"
	companyIDKey    contextKey = "companyID"
	roleIDKey       contextKey = "roleID"
	permissionIDKey contextKey = "permissionID"
)

func GetLanguage(ctx context.Context) string {
	if lang, ok := ctx.Value(LanguageKey).(string); ok {
		return lang
	}
	return language.English // Default to English
}

func SetLanguage(ctx context.Context, lang string) context.Context {
	return context.WithValue(ctx, LanguageKey, lang)
}

func LanguageMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lang := r.Header.Get("Accept-Language")
			if lang == "" {
				lang = language.English
			}

			ctx := SetLanguage(r.Context(), lang)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
