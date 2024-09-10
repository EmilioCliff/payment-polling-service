package services

import (
	"time"

	"github.com/google/uuid"
)

type RegisterUserRequest struct {
	FullName       string `json:"full_name" binding:"required"`
	Email          string `json:"email" binding:"required"`
	Password       string `json:"password" binding:"required"`
	PaydUsername   string `json:"payd_username" binding:"required"`
	UsernameApiKey string `json:"payd_username_api_key" binding:"required"`
	PasswordApiKey string `json:"payd_password_api_key" binding:"required"`
}

type RegisterUserResponse struct {
	FullName   string    `json:"full_name"`
	Email      string    `json:"email"`
	CreatedAt  time.Time `json:"created_at"`
	Message    string    `json:"message,omitempty"`
	StatusCode int       `json:"status_code,omitempty"`
}

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginUserResponse struct {
	AccessToken string    `json:"access_token"`
	FullName    string    `json:"full_name"`
	Email       string    `json:"email"`
	CreatedAt   time.Time `json:"created_at"`
	Message     string    `json:"message,omitempty"`
	StatusCode  int       `json:"status_code,omitempty"`
}

type InitiatePaymentRequest struct {
	UserID      int64  `json:"user_id" binding:"required"`
	Action      string `json:"action" binding:"required"`
	Amount      int64  `json:"amount" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	NetworkCode string `json:"network_code" binding:"required"`
	Naration    string `json:"naration" binding:"required"`
}

type InitiatePaymentResponse struct {
	TransactionID string `json:"transaction_id"`
	PaymentStatus bool   `json:"payment_status"`
	Action        string `json:"action"`
	Message       string `json:"message,omitempty"`
	StatusCode    int    `json:"status_code,omitempty"`
}

type PollingTransactionRequest struct {
	TransactionId string `uri:"id" binding:"required"`
}

type PollingTransactionResponse struct {
	TransactionID      uuid.UUID `json:"transaction_id"`
	PaydTransactionRef string    `json:"payd_transaction_ref"`
	Action             string    `json:"action"`
	Amount             int64     `json:"amount"`
	PhoneNumber        string    `json:"phone_number"`
	NetworkCode        string    `json:"network_code"`
	Naration           string    `json:"naration"`
	PaymentStatus      bool      `json:"payment_status"`
	PaydUsername       string    `json:"payd_username"`
	PaydUsernameApiKey string    `json:"payd_username_api_key"`
	PaydPasswordApiKey string    `json:"payd_password_api_key"`
	Message            string    `json:"message,omitempty"`
	StatusCode         int       `json:"status_code,omitempty"`
}
