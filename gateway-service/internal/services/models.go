package services

import (
	"time"

	"github.com/google/uuid"
)

type RegisterUserRequest struct {
	FullName       string `binding:"required" json:"full_name"`
	Email          string `binding:"required" json:"email"`
	Password       string `binding:"required" json:"password"`
	PaydUsername   string `binding:"required" json:"payd_username"`
	PaydAccountID  string `binding:"required" json:"payd_account_id"`
	UsernameApiKey string `binding:"required" json:"username_api_key"`
	PasswordApiKey string `binding:"required" json:"password_api_key"`
}

type RegisterUserResponse struct {
	FullName   string    `json:"full_name,omitempty"`
	Email      string    `json:"email,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	Message    string    `json:"message,omitempty"`
	StatusCode int       `json:"status_code,omitempty"`
}

type LoginUserRequest struct {
	Email    string `binding:"required" json:"email"`
	Password string `binding:"required" json:"password"`
}

type LoginUserResponse struct {
	AccessToken  string    `json:"access_token,omitempty"`
	FullName     string    `json:"full_name,omitempty"`
	Email        string    `json:"email,omitempty"`
	ExpirationAt time.Time `json:"expiration_at,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	Message      string    `json:"message,omitempty"`
	StatusCode   int       `json:"status_code,omitempty"`
}

type InitiatePaymentRequest struct {
	UserID      int64  `binding:"required"                          json:"user_id"`
	Action      string `binding:"required,oneof=withdrawal payment" json:"action"`
	Amount      int64  `binding:"required"                          json:"amount"`
	PhoneNumber string `binding:"required"                          json:"phone_number"`
	NetworkCode string `binding:"required,oneof=63902 63903"        json:"network_code"`
	Naration    string `binding:"required"                          json:"naration"`
}

type InitiatePaymentResponse struct {
	TransactionID string `json:"transaction_id,omitempty"`
	PaymentStatus bool   `json:"payment_status,omitempty"`
	Action        string `json:"action,omitempty"`
	Message       string `json:"message,omitempty"`
	StatusCode    int    `json:"status_code,omitempty"`
}

type PollingTransactionRequest struct {
	TransactionId string `binding:"required" uri:"id"`
}

type PollingTransactionResponse struct {
	TransactionID      uuid.UUID `json:"transaction_id,omitempty"`
	PaydTransactionRef string    `json:"payd_transaction_ref,omitempty"`
	Action             string    `json:"action,omitempty"`
	Amount             int64     `json:"amount,omitempty"`
	PhoneNumber        string    `json:"phone_number,omitempty"`
	NetworkCode        string    `json:"network_code,omitempty"`
	Naration           string    `json:"naration,omitempty"`
	PaymentStatus      bool      `json:"payment_status,omitempty"`
	Message            string    `json:"message,omitempty"`
	StatusCode         int       `json:"status_code,omitempty"`
}
