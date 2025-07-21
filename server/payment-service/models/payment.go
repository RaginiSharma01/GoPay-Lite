package models

import (
	"time"
)

// PaymentRequest represents the incoming payment request payload
// @swagger:model PaymentRequest
type PaymentRequest struct {
	// Amount in the specified currency (must be positive)
	// required: true
	// minimum: 0.01
	Amount float64 `json:"amount" validate:"required,gt=0"`

	// Source account ID
	// required: true
	// example: acc_123456789
	FromAccount string `json:"from_account" validate:"required"`

	// Destination account ID
	// required: true
	// example: acc_987654321
	ToAccount string `json:"to_account" validate:"required"`

	// Currency code (ISO 4217)
	// default: "INR"
	// example: INR
	Currency string `json:"currency,omitempty" validate:"omitempty,len=3"`
}

// PaymentResponse represents the API response for a successful payment
// @swagger:model PaymentResponse
type PaymentResponse struct {
	// Payment ID
	// example: 42
	ID int `json:"id"`

	// Razorpay order ID
	// example: order_123456789
	RazorpayOrderID string `json:"razorpay_order_id,omitempty"`

	// Payment status
	// example: created
	Status string `json:"status"`

	// Payment amount
	// example: 100.50
	Amount float64 `json:"amount"`

	// Currency code
	// example: INR
	Currency string `json:"currency,omitempty"`

	// Timestamp of creation
	// example: 2023-05-15T14:30:45Z
	CreatedAt time.Time `json:"created_at"`

	// Timestamp of last update
	// example: 2023-05-15T14:30:45Z
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// Payment represents the payment record in database
type Payment struct {
	ID              int        `json:"id" db:"id"`
	UserID          int        `json:"user_id" db:"user_id"`
	Amount          float64    `json:"amount" db:"amount"`
	Currency        string     `json:"currency" db:"currency"`
	FromAccount     string     `json:"from_account" db:"from_account"`
	ToAccount       string     `json:"to_account" db:"to_account"`
	RazorpayOrderID string     `json:"razorpay_order_id" db:"razorpay_order_id"`
	Status          string     `json:"status" db:"status"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty" db:"updated_at"`
	Description     *string    `json:"description,omitempty" db:"description"`
}

// ErrorResponse represents standard API error response
// @swagger:model ErrorResponse
type ErrorResponse struct {
	// Error type
	// example: validation_error
	Error string `json:"error"`

	// Human-readable message
	// example: Amount must be positive
	Message string `json:"message"`

	// Optional field-specific errors
	Errors map[string]string `json:"errors,omitempty"`
}

// PaymentStatus represents possible payment states
type PaymentStatus string

const (
	PaymentStatusCreated   PaymentStatus = "created"
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)
