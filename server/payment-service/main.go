package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/RaginiSharma01/gopay-lite/payment-service/db"
	"github.com/RaginiSharma01/gopay-lite/payment-service/handlers"
	"github.com/RaginiSharma01/gopay-lite/payment-service/middleware"
	"github.com/RaginiSharma01/gopay-lite/payment-service/razorpay"

	_ "github.com/RaginiSharma01/gopay-lite/payment-service/docs" // Swagger generated docs
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title GoPay Payment Service API
// @version 1.0
// @description This service handles payment transactions.
// @host localhost:8084
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// Load environment variables
	if os.Getenv("RUNNING_IN_DOCKER") != "true" {
		err := godotenv.Load("../.env")
		if err != nil {
			log.Println("Warning: .env file not found or failed to load")
		}
	}

	// Initialize database
	if err := db.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize Razorpay client
	razorpay.Init()
	razorpayKey := os.Getenv("RAZORPAY_KEY")
	razorpaySecret := os.Getenv("RAZORPAY_SECRET")

	if razorpayKey == "" || razorpaySecret == "" {
		log.Fatal("Razorpay credentials not set in environment")
	}

	handlers.InitRazorpayClient(razorpayKey, razorpaySecret)

	// Create router
	r := mux.NewRouter()

	// Welcome endpoint
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to GoPay Payment Service"))
	}).Methods("GET")

	// Middlewares
	r.Use(loggingMiddleware)
	r.Use(corsMiddleware) // 👈 NEW: Add CORS middleware here
	r.Use(middleware.ContentTypeJSON)

	// Swagger docs
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Health check
	r.HandleFunc("/health", healthCheck).Methods("GET")

	// Protected routes
	api := r.PathPrefix("/api/v1").Subrouter()
	api.Use(middleware.JWTAuth)
	api.HandleFunc("/pay", handlers.HandlePayment).Methods("POST")

	// Server setup
	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Payment Service running on http://localhost:%s", port)
		log.Printf("Swagger UI available at http://localhost:%s/swagger/index.html", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-done
	log.Println("Shutting down Payment Service...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown error: %v", err)
	}

	log.Println("Stopped gracefully")
}

// Health check endpoint
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// Logging middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// CORS middleware for allowing frontend requests
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "http://localhost:3000" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
