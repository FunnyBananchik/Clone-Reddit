package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type ctxKey string

const UserContextKey ctxKey = "userKey"

func (h *AuthHandler) AuthMiddlewear(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			errorJSON(w, 401, "unauthorized")
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			errorJSON(w, 401, "bad token")
			return
		}
		token := parts[1]
		claims := &Claims{}
		t, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
			return h.secret, nil
		}, jwt.WithValidMethods([]string{"HS256"}), jwt.WithExpirationRequired(), jwt.WithIssuedAt())
		if err != nil || !t.Valid {
			errorJSON(w, 401, "bad token")
			return
		}
		ctx := context.WithValue(r.Context(), UserContextKey, claims.User)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
