package middleware

import (
	"context"
	"delivery_app/utils"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

var BlacklistedTokens = make(map[string]bool)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if BlacklistedTokens[tokenString] {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		token, err := utils.VerifyToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		id, ok := claims["id"].(float64)
		if !ok {
			http.Error(w, "Invalid user ID in token claims", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, int(id))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
