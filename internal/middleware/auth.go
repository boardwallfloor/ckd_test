package middleware

import (
	"boardwallfloor/ckd/internal/service"
	"context"
	"net/http"
	"strings"
)

type contextKey string

const AuthClaimsKey contextKey = "authClaims"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"success": false, "message": "Authorization header required"}`, http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, `{"success": false, "message": "Invalid authorization header format"}`, http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		claims, err := service.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, `{"success": false, "message": "Invalid token"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), AuthClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
