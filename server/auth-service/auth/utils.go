package auth

import (
	"log"
	"os"
)

var jwtSecret []byte

func init() {
	// Try to load from environment
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))

	// Fallback for development (remove in production)
	if len(jwtSecret) == 0 {
		log.Println("WARNING: Using development JWT secret - configure JWT_SECRET in production")
		jwtSecret = []byte("dev-secret-please-change-me")
	}
}
