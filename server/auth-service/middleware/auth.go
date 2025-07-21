package middleware

import (
	"context"
	"net/http"
	"strings"

	"gopay-lite/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	emailKey contextKey = "email"
	// Consider adding more context keys for other claims like userID, roles, etc.
)

// VerifyJWT is a middleware that checks for valid JWT token
func VerifyJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Extract and validate Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Authorization header must start with 'Bearer '", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 2. Parse and validate JWT
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(config.GetJWTSecret()), nil
		})

		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// 3. Extract and validate claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims format", http.StatusUnauthorized)
			return
		}

		email, ok := claims["email"].(string)
		if !ok || email == "" {
			http.Error(w, "Email claim missing or invalid", http.StatusUnauthorized)
			return
		}

		// 4. Add claims to context
		ctx := context.WithValue(r.Context(), emailKey, email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetEmailFromContext safely extracts email from the request context
func GetEmailFromContext(r *http.Request) (string, bool) {
	email, ok := r.Context().Value(emailKey).(string)
	return email, ok
}
