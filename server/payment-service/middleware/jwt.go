package middleware

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	EmailKey  contextKey = "email"
	UserIDKey contextKey = "userID"
)

// Exported for access in handlers
var (
	UserIDContextKey = UserIDKey
)

// JWTAuth validates JWT tokens and sets user info in context
func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Println("Missing Authorization header")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			http.Error(w, "Bearer token required", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		}, jwt.WithValidMethods([]string{"HS256"}))

		if err != nil || !token.Valid {
			log.Printf("Token parsing error: %v", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			email, _ := claims["email"].(string)
			uidFloat, ok := claims["user_id"].(float64)
			if !ok {
				http.Error(w, "Invalid token claims: user_id missing", http.StatusUnauthorized)
				return
			}
			userID := int(uidFloat)

			ctx := context.WithValue(r.Context(), EmailKey, email)
			ctx = context.WithValue(ctx, UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		http.Error(w, "Failed to parse claims", http.StatusUnauthorized)
	})
}

// ContentTypeJSON sets response content type to application/json
func ContentTypeJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
