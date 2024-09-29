package repository

import (
	"context"
	"time"

	"github.com/EmilioCliff/payment-polling-app/payment-service/pkg"
	"github.com/google/uuid"
)

type Transaction struct {
	TransactionID      uuid.UUID `json:"transaction_id"`
	PaydTransactionRef string    `json:"payd_transaction_ref"`
	Message            string    `json:"message"`
	UserID             int64     `json:"user_id"`
	Action             string    `json:"action"`
	Amount             int32     `json:"amount"`
	PhoneNumber        string    `json:"phone_number"`
	NetworkCode        string    `json:"network_code"`
	Narration          string    `json:"narration"`
	Status             bool      `json:"status"`
	UpdatedAt          time.Time `json:"updated_at"`
	CreatedAt          time.Time `json:"created_at"`
}

type TransactionUpdate struct {
	PaydTransactionRef string `json:"payd_transaction_ref"`
	Status             bool   `json:"status"`
	Message            string `json:"message"`
}

func (t *Transaction) Validate() error {
	if t.TransactionID == uuid.Nil {
		return pkg.Errorf(pkg.INVALID_ERROR, "transaction_id is required")
	}

	if t.UserID == 0 {
		return pkg.Errorf(pkg.INVALID_ERROR, "user_id is required")
	}

	if t.Action == "" || t.Action != "withdrawal" && t.Action != "payment" {
		return pkg.Errorf(pkg.INVALID_ERROR, "action is invalid")
	}

	if t.Amount == 0 {
		return pkg.Errorf(pkg.INVALID_ERROR, "amount is required")
	}

	if t.PhoneNumber == "" {
		return pkg.Errorf(pkg.INVALID_ERROR, "phone_number is required")
	}

	if t.NetworkCode == "" || t.NetworkCode != "63902" && t.NetworkCode != "63903" {
		return pkg.Errorf(pkg.INVALID_ERROR, "network_node is invalid")
	}

	if t.Narration == "" {
		return pkg.Errorf(pkg.INVALID_ERROR, "narration is required")
	}

	return nil
}

type TransactionRepository interface {
	CreateTransaction(context.Context, Transaction) (*Transaction, error)
	PollingTransaction(context.Context, uuid.UUID) (*Transaction, error)
	UpdateTransaction(context.Context, uuid.UUID, TransactionUpdate) (*Transaction, error)
}
