package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"gopay-lite/db"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// ========== Data Models ==========

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
}

type MeResponse struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
}

// ========== Register ==========

// Register a new user
//
// @Summary Register a user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param user body User true "User data"
// @Success 201 {object} AuthResponse
// @Failure 400 {object} AuthResponse
// @Failure 409 {object} AuthResponse
// @Failure 500 {object} AuthResponse
// @Router /api/v1/register [post]
func Register(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		sendErrorResponse(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	u.Name = strings.TrimSpace(u.Name)
	u.Email = strings.TrimSpace(u.Email)
	u.Password = strings.TrimSpace(u.Password)

	if u.Name == "" || u.Email == "" || u.Password == "" {
		sendErrorResponse(w, "Name, email, and password are required", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		sendErrorResponse(w, "Error processing password", http.StatusInternalServerError)
		return
	}

	// Insert user and get ID
	var userID int
	err = db.DB.QueryRow(
		"INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id",
		u.Name, u.Email, string(hashedPassword),
	).Scan(&userID)

	if err != nil {
		sendErrorResponse(w, "Email already registered", http.StatusConflict)
		return
	}

	// Generate JWT
	token, err := generateToken(u.Email, userID)
	if err != nil {
		sendErrorResponse(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	// Respond
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(AuthResponse{
		Message: "User registered successfully",
		Token:   token,
	})
}

// ========== Login ==========

// Login a user and return JWT
//
// @Summary Login a user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param user body User true "User credentials"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} AuthResponse
// @Failure 401 {object} AuthResponse
// @Failure 500 {object} AuthResponse
// @Router /api/v1/login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		sendErrorResponse(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	u.Email = strings.TrimSpace(u.Email)
	u.Password = strings.TrimSpace(u.Password)

	if u.Email == "" || u.Password == "" {
		sendErrorResponse(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	var userID int
	var hashedPassword string

	err := db.DB.QueryRow("SELECT id, password FROM users WHERE email = $1", u.Email).Scan(&userID, &hashedPassword)
	if err != nil {
		sendErrorResponse(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(u.Password)); err != nil {
		sendErrorResponse(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token, err := generateToken(u.Email, userID)
	if err != nil {
		sendErrorResponse(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(AuthResponse{
		Message: "Login successful",
		Token:   token,
	})
}

// ========== Me (JWT Protected) ==========

// Me handler returns user info from token
//
// @Summary Get user info
// @Description Returns email and user ID from JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} MeResponse
// @Failure 401 {object} AuthResponse
// @Failure 500 {object} AuthResponse
// @Router /api/v1/me [get]
func Me(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		sendErrorResponse(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	tokenString := extractToken(authHeader)
	if tokenString == "" {
		sendErrorResponse(w, "Invalid Authorization header format", http.StatusUnauthorized)
		return
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		sendErrorResponse(w, "JWT_SECRET not set", http.StatusInternalServerError)
		return
	}

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		sendErrorResponse(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	email, _ := claims["email"].(string)
	userID, _ := claims["user_id"].(float64)

	json.NewEncoder(w).Encode(MeResponse{
		UserID: int(userID),
		Email:  email,
	})
}

// ========== JWT Utility ==========

func generateToken(email string, userID int) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET not set")
	}

	claims := jwt.MapClaims{
		"email":   email,
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(),
		"iss":     "gopay-lite",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func extractToken(authHeader string) string {
	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
		return parts[1]
	}
	return ""
}

// ========== Helper ==========

func sendErrorResponse(w http.ResponseWriter, msg string, code int) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(AuthResponse{Message: msg})
}
