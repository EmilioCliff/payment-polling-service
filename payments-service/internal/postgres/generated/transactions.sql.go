// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: transactions.sql

package generated

import (
	"context"

	"github.com/google/uuid"
)

const createTransaction = `-- name: CreateTransaction :one
INSERT INTO transactions (
    transaction_id, payd_transaction_ref,user_id, message, action, amount, phone_number, network_node, narration
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING transaction_id, payd_transaction_ref, user_id, action, amount, phone_number, network_node, narration, status, updated_at, created_at, message
`

type CreateTransactionParams struct {
	TransactionID      uuid.UUID `json:"transaction_id"`
	PaydTransactionRef string    `json:"payd_transaction_ref"`
	UserID             int64     `json:"user_id"`
	Message            string    `json:"message"`
	Action             string    `json:"action"`
	Amount             int32     `json:"amount"`
	PhoneNumber        string    `json:"phone_number"`
	NetworkNode        string    `json:"network_node"`
	Narration          string    `json:"narration"`
}

func (q *Queries) CreateTransaction(ctx context.Context, arg CreateTransactionParams) (Transaction, error) {
	row := q.db.QueryRow(ctx, createTransaction,
		arg.TransactionID,
		arg.PaydTransactionRef,
		arg.UserID,
		arg.Message,
		arg.Action,
		arg.Amount,
		arg.PhoneNumber,
		arg.NetworkNode,
		arg.Narration,
	)
	var i Transaction
	err := row.Scan(
		&i.TransactionID,
		&i.PaydTransactionRef,
		&i.UserID,
		&i.Action,
		&i.Amount,
		&i.PhoneNumber,
		&i.NetworkNode,
		&i.Narration,
		&i.Status,
		&i.UpdatedAt,
		&i.CreatedAt,
		&i.Message,
	)
	return i, err
}

const getTransaction = `-- name: GetTransaction :one
SELECT transaction_id, payd_transaction_ref, user_id, action, amount, phone_number, network_node, narration, status, updated_at, created_at, message FROM transactions
WHERE transaction_id = $1
`

func (q *Queries) GetTransaction(ctx context.Context, transactionID uuid.UUID) (Transaction, error) {
	row := q.db.QueryRow(ctx, getTransaction, transactionID)
	var i Transaction
	err := row.Scan(
		&i.TransactionID,
		&i.PaydTransactionRef,
		&i.UserID,
		&i.Action,
		&i.Amount,
		&i.PhoneNumber,
		&i.NetworkNode,
		&i.Narration,
		&i.Status,
		&i.UpdatedAt,
		&i.CreatedAt,
		&i.Message,
	)
	return i, err
}

const updateTransaction = `-- name: UpdateTransaction :one
UPDATE transactions
SET status = $2,
    updated_at = now(),
    payd_transaction_ref = $3,
    message = $4
WHERE transaction_id = $1
RETURNING transaction_id, payd_transaction_ref, user_id, action, amount, phone_number, network_node, narration, status, updated_at, created_at, message
`

type UpdateTransactionParams struct {
	TransactionID      uuid.UUID `json:"transaction_id"`
	Status             bool      `json:"status"`
	PaydTransactionRef string    `json:"payd_transaction_ref"`
	Message            string    `json:"message"`
}

func (q *Queries) UpdateTransaction(ctx context.Context, arg UpdateTransactionParams) (Transaction, error) {
	row := q.db.QueryRow(ctx, updateTransaction,
		arg.TransactionID,
		arg.Status,
		arg.PaydTransactionRef,
		arg.Message,
	)
	var i Transaction
	err := row.Scan(
		&i.TransactionID,
		&i.PaydTransactionRef,
		&i.UserID,
		&i.Action,
		&i.Amount,
		&i.PhoneNumber,
		&i.NetworkNode,
		&i.Narration,
		&i.Status,
		&i.UpdatedAt,
		&i.CreatedAt,
		&i.Message,
	)
	return i, err
}
