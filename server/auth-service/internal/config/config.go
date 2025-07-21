package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables.
// In dev: loads from ../.env using godotenv.
// In Docker: relies on env vars passed at runtime.
func LoadEnv() {
	if os.Getenv("RUNNING_IN_DOCKER") != "true" {
		err := godotenv.Load("../.env")
		if err != nil {
			log.Println("⚠️  Warning: Could not load .env file (for dev only)")
		} else {
			log.Println("✅ .env file loaded successfully for dev")
		}
	} else {
		log.Println("🚀 Running in Docker – using passed env variables")
	}

	// 🔍 Debug env values (you can remove this after testing)
	log.Println("RUNNING_IN_DOCKER:", os.Getenv("RUNNING_IN_DOCKER"))
	log.Println("DB_HOST:", os.Getenv("DB_HOST"))
	log.Println("DB_PORT:", os.Getenv("DB_PORT"))
	log.Println("DB_USER:", os.Getenv("DB_USER"))
	log.Println("DB_NAME:", os.Getenv("DB_NAME"))
	log.Println("DB_PASSWORD:", os.Getenv("DB_PASSWORD"))
}

// GetJWTSecret retrieves the JWT secret from environment variables.
func GetJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Println("⚠️  Warning: JWT_SECRET is not set")
	}
	return secret
}
