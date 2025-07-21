package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/RaginiSharma01/gopay-lite/payment-service/db"
	"github.com/RaginiSharma01/gopay-lite/payment-service/middleware"
	"github.com/RaginiSharma01/gopay-lite/payment-service/models"
	razorpay "github.com/razorpay/razorpay-go"
)

var razorpayClient *razorpay.Client

// Payment represents the payment model
type Payment struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	Amount          float64   `json:"amount"`
	Currency        string    `json:"currency"`
	FromAccount     string    `json:"from_account"`
	ToAccount       string    `json:"to_account"`
	RazorpayOrderID string    `json:"razorpay_order_id"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}

// InitRazorpayClient initializes the Razorpay client
func InitRazorpayClient(key, secret string) {
	razorpayClient = razorpay.NewClient(key, secret)
}

// HandlePayment handles the payment request
// @Summary Process payment
// @Description Process a payment transaction
// @Tags payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payment body models.PaymentRequest true "Payment details"
// @Success 201 {object} models.PaymentResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/pay [post]
func HandlePayment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Decode request
	var req models.PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "Invalid request",
			Message: "Failed to parse request body",
		})
		return
	}
	log.Printf("Payment request received: %+v", req)

	// Validate request
	if req.Amount <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "Invalid amount",
			Message: "Amount must be positive",
		})
		return
	}

	if req.FromAccount == "" || req.ToAccount == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "Missing accounts",
			Message: "Both from_account and to_account must be specified",
		})
		return
	}

	// Set default currency if not provided
	if req.Currency == "" {
		req.Currency = "INR"
	}

	// Get user from context
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(int)

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "Unauthorized",
			Message: "Invalid user context",
		})
		return
	}

	// Begin database transaction
	tx, err := db.DB.BeginTx(r.Context(), nil)
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "Database error",
			Message: "Failed to start transaction",
		})
		return
	}
	defer tx.Rollback()

	// Create Razorpay order
	orderData := map[string]interface{}{
		"amount":   int(req.Amount * 100), // Convert to paise
		"currency": req.Currency,
		"receipt":  fmt.Sprintf("order_%d_%d", userID, time.Now().Unix()),
		"notes": map[string]interface{}{
			"from_account": req.FromAccount,
			"to_account":   req.ToAccount,
		},
	}

	order, err := razorpayClient.Order.Create(orderData, nil)
	if err != nil {
		log.Printf("Razorpay order creation failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "Payment failed",
			Message: "Could not create payment order",
		})
		return
	}

	// Prepare payment record
	payment := Payment{
		UserID:          userID,
		Amount:          req.Amount,
		Currency:        req.Currency,
		FromAccount:     req.FromAccount,
		ToAccount:       req.ToAccount,
		RazorpayOrderID: order["id"].(string),
		Status:          "created",
		CreatedAt:       time.Now().UTC(),
	}

	// Store payment in database
	query := `INSERT INTO payments 
		(user_id, amount, currency, from_account, to_account, razorpay_order_id, status, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		RETURNING id`

	err = tx.QueryRowContext(r.Context(),
		query,
		payment.UserID,
		payment.Amount,
		payment.Currency,
		payment.FromAccount,
		payment.ToAccount,
		payment.RazorpayOrderID,
		payment.Status,
		payment.CreatedAt,
	).Scan(&payment.ID)

	if err != nil {
		log.Printf("Database error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "Payment processing failed",
			Message: "Could not save payment record",
		})
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Transaction commit failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "Transaction error",
			Message: "Failed to complete payment",
		})
		return
	}

	// Return success response
	response := models.PaymentResponse{
		ID:              payment.ID,
		RazorpayOrderID: payment.RazorpayOrderID,
		Status:          payment.Status,
		Amount:          payment.Amount,
		Currency:        payment.Currency,
		CreatedAt:       payment.CreatedAt,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
