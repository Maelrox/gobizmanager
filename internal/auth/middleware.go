package auth

import (
	"context"
	"net/http"
	"strings"

	appcontext "gobizmanager/pkg/context"
	"gobizmanager/pkg/language"
	"gobizmanager/pkg/utils"
)

type contextKey string

const UserIDKey contextKey = "userID"

func Middleware(jwtManager *JWTManager, msgStore *language.MessageStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lang := appcontext.GetLanguage(r.Context())

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.JSONError(w, http.StatusUnauthorized, msgStore.GetMessage(lang, "auth.authorization_required"))
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				utils.JSONError(w, http.StatusUnauthorized, msgStore.GetMessage(lang, "auth.invalid_authorization_format"))
				return
			}

			tokenString := parts[1]
			claims, err := jwtManager.VerifyToken(tokenString)
			if err != nil {
				if err == ErrExpiredToken {
					utils.JSONError(w, http.StatusUnauthorized, msgStore.GetMessage(lang, "auth.token_expired"))
				} else {
					utils.JSONError(w, http.StatusUnauthorized, msgStore.GetMessage(lang, "auth.invalid_token"))
				}
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(UserIDKey).(int64)
	return userID, ok
}
