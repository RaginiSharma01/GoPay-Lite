package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/RaginiSharma01/gopay-lite/api-gateway/routes"
	"github.com/gorilla/mux"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set your frontend origin explicitly
		origin := r.Header.Get("Origin")
		allowedOrigins := map[string]bool{
			"http://localhost:3000": true,
			// Add other allowed origins here
		}

		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
func main() {
	// Load service URLs from environment
	authURL := getEnv("AUTH_SERVICE_URL", "http://localhost:8083")
	paymentURL := getEnv("PAYMENT_SERVICE_URL", "http://localhost:8084")

	r := mux.NewRouter()
	r.Use(enableCORS)
	r.Use(loggingMiddleware)

	// Route configurations
	r.PathPrefix("/api/v1/auth/").Handler(
		routes.NewReverseProxy(authURL, "/api/v1/auth", "/api/v1"),
	)

	r.PathPrefix("/api/v1/pay").Handler(
		routes.NewReverseProxy(paymentURL, "/api/v1", "/api/v1"),
	)

	// Health check endpoint
	r.HandleFunc("/health", healthCheck).Methods("GET")

	// Server configuration
	server := &http.Server{
		Addr:         ":" + getEnv("PORT", "8080"),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("🚀 API Gateway running on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("💥 Server failed: %v", err)
	}
}

// Helper functions
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}
