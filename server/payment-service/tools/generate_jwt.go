package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env if running locally
	err := godotenv.Load("../.env") // Adjust path if needed
	if err != nil {
		log.Println("Warning: .env file not found, trying environment variable directly")
	}

	// Get secret key from environment
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is not set in environment")
	}

	// Prepare claims that match middleware expectations
	claims := jwt.MapClaims{
		"email": "ragini@example.com",                  // Required by middleware
		"exp":   time.Now().Add(24 * time.Hour).Unix(), // Expiration
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token using the secret
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Fatalf("Error signing token: %v", err)
	}

	// Print token to terminal
	fmt.Println("\nGenerated JWT Token:")
	fmt.Println(signedToken)
}
