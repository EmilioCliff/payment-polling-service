package postgres

import (
	"context"
	"database/sql"

	"github.com/EmilioCliff/payment-polling-app/payment-service/internal/postgres/generated"
	"github.com/EmilioCliff/payment-polling-app/payment-service/pkg"
	"github.com/google/uuid"
)

type InitiatePaymentRequest struct {
	TransactionID      uuid.UUID `json:"transaction_id"`
	UserID             int64     `json:"user_id"`
	Action             string    `json:"action"`
	Amount             int64     `json:"amount"`
	PhoneNumber        string    `json:"phone_number"`
	NetworkCode        string    `json:"network_code"`
	Naration           string    `json:"naration"`
	PaydTransactionRef string    `json:"payd_transaction_ref"`
}

type InitiatePaymentResponse struct {
	TransactionID uuid.UUID `json:"transaction_id"`
	PaymentStatus bool      `json:"payment_status"`
	Action        string    `json:"action"`
	Amount        int32     `json:"amount"`
}

func (store *Store) CreateTransaction(ctx context.Context, req InitiatePaymentRequest) (*InitiatePaymentResponse, *pkg.Error) {
	// validate the request

	transaction, err := store.Queries.CreateTransaction(ctx, generated.CreateTransactionParams{
		TransactionID:      req.TransactionID,
		PaydTransactionRef: req.PaydTransactionRef,
		UserID:             req.UserID,
		Action:             req.Action,
		Amount:             int32(req.Amount),
		PhoneNumber:        req.PhoneNumber,
		NetworkNode:        req.NetworkCode,
		Narration:          req.Naration,
	})
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to create transaction: %v", err)
	}

	return &InitiatePaymentResponse{
		TransactionID: transaction.TransactionID,
		PaymentStatus: transaction.Status,
		Action:        transaction.Action,
		Amount:        transaction.Amount,
	}, nil
}

type PollingTransactionRequest struct {
	TransactionId string `json:"transaction_id"`
}

type PollingTransactionResponse struct {
	TransactionID      uuid.UUID `json:"transaction_id"`
	PaydTransactionRef string    `json:"payd_transaction_ref"`
	Action             string    `json:"action"`
	Amount             int32     `json:"amount"`
	PhoneNumber        string    `json:"phone_number"`
	NetworkCode        string    `json:"network_code"`
	Naration           string    `json:"naration"`
	PaymentStatus      bool      `json:"payment_status"`
}

func (store *Store) PollingTransaction(ctx context.Context, req PollingTransactionRequest) (*PollingTransactionResponse, *pkg.Error) {
	transactionUUID, err := uuid.Parse(req.TransactionId)
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "invalid uuid: %v", err)
	}

	transaction, err := store.Queries.GetTransaction(ctx, transactionUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "transaction does not exist: %v", err)
		}
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error getting transaction: %v", err)
	}

	return &PollingTransactionResponse{
		TransactionID:      transaction.TransactionID,
		PaydTransactionRef: transaction.PaydTransactionRef,
		Action:             transaction.Action,
		Amount:             transaction.Amount,
		PhoneNumber:        transaction.PhoneNumber,
		NetworkCode:        transaction.NetworkNode,
		Naration:           transaction.Narration,
		PaymentStatus:      transaction.Status,
	}, nil
}
