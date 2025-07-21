package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gopay-lite/auth"
	"gopay-lite/db"
	_ "gopay-lite/docs" // Swagger generated docs
	"gopay-lite/internal/config"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

// @title GoPay-Lite Auth API
// @version 1.0
// @description Secure Authentication Service for GoPay-Lite
// @contact.name API Support
// @contact.email support@gopay-lite.com
// @license.name MIT
// @host localhost:8083
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	//  Load env using internal/config
	config.LoadEnv()

	// Now your DB_*, JWT_SECRET etc. should be available to os.Getenv()
	db.Init()

	// Router setup
	r := mux.NewRouter()

	// Middleware
	r.Use(loggingMiddleware)
	r.Use(mux.CORSMethodMiddleware(r))

	// Swagger docs
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // Adjust if needed
	))

	// Base handlers
	r.HandleFunc("/health", healthCheck).Methods("GET")
	r.HandleFunc("/", rootHandler).Methods("GET")

	// Auth API group with versioned prefix
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/register", auth.Register).Methods("POST")
	api.HandleFunc("/login", auth.Login).Methods("POST")
	api.HandleFunc("/me", auth.Me).Methods("GET")

	// Server setup
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	corsOpts := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}), // or "*" for any origin (not for production)
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      corsOpts(r), // Wrap the router with CORS
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Logging
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Starting Auth Service", zap.String("port", port), zap.Time("started_at", time.Now()))

	go func() {
		log.Printf("Auth Service running at http://localhost:%s", port)
		log.Printf("Swagger: http://localhost:%s/swagger/index.html", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-done
	log.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown error: %v", err)
	}
	log.Println("Stopped gracefully")
}

//============ Handlers ============

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GoPay-Lite Auth Service - See /swagger for docs"))
}

// ============ Logging Middleware ============

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("[REQUEST] %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("[COMPLETED] %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	})
}
