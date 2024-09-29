package postgres

import (
	"context"
	"database/sql"

	"github.com/EmilioCliff/payment-polling-app/payment-service/internal/postgres/generated"
	"github.com/EmilioCliff/payment-polling-app/payment-service/internal/repository"
	"github.com/EmilioCliff/payment-polling-app/payment-service/pkg"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

var _ repository.TransactionRepository = (*TransactionRepository)(nil)

type TransactionRepository struct {
	db      *Store
	queries generated.Querier
}

func NewTransactionService(db *Store) *TransactionRepository {
	queries := generated.New(db.conn)

	return &TransactionRepository{
		db:      db,
		queries: queries,
	}
}

func (t *TransactionRepository) CreateTransaction(ctx context.Context, transaction repository.Transaction) (*repository.Transaction, error) {
	err := transaction.Validate()
	if err != nil {
		return nil, err
	}

	createdTransaction, err := t.queries.CreateTransaction(ctx, generated.CreateTransactionParams{
		TransactionID:      transaction.TransactionID,
		PaydTransactionRef: transaction.PaydTransactionRef,
		Message:            transaction.Message,
		UserID:             transaction.UserID,
		Action:             transaction.Action,
		Amount:             transaction.Amount,
		PhoneNumber:        transaction.PhoneNumber,
		NetworkNode:        transaction.NetworkCode,
		Narration:          transaction.Narration,
	})
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				return nil, pkg.Errorf(pkg.ALREADY_EXISTS_ERROR, "transaction already exists")
			}
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to create transaction")
	}

	return &repository.Transaction{
		TransactionID:      createdTransaction.TransactionID,
		PaydTransactionRef: createdTransaction.PaydTransactionRef,
		UserID:             createdTransaction.UserID,
		Action:             createdTransaction.Action,
		Amount:             createdTransaction.Amount,
		PhoneNumber:        createdTransaction.PhoneNumber,
		NetworkCode:        createdTransaction.NetworkNode,
		Narration:          createdTransaction.Narration,
		Status:             createdTransaction.Status,
		UpdatedAt:          createdTransaction.UpdatedAt,
		CreatedAt:          createdTransaction.CreatedAt,
	}, nil
}

func (t *TransactionRepository) PollingTransaction(ctx context.Context, id uuid.UUID) (*repository.Transaction, error) {
	transaction, err := t.queries.GetTransaction(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "transaction does not exist")
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error getting transaction")
	}

	return &repository.Transaction{
		TransactionID:      transaction.TransactionID,
		PaydTransactionRef: transaction.PaydTransactionRef,
		Message:            transaction.Message,
		UserID:             transaction.UserID,
		Action:             transaction.Action,
		Amount:             transaction.Amount,
		PhoneNumber:        transaction.PhoneNumber,
		NetworkCode:        transaction.NetworkNode,
		Narration:          transaction.Narration,
		Status:             transaction.Status,
		UpdatedAt:          transaction.UpdatedAt,
		CreatedAt:          transaction.CreatedAt,
	}, nil
}

func (t *TransactionRepository) UpdateTransaction(
	ctx context.Context,
	id uuid.UUID,
	update repository.TransactionUpdate,
) (*repository.Transaction, error) {
	transaction, err := t.queries.UpdateTransaction(ctx, generated.UpdateTransactionParams{
		TransactionID:      id,
		Status:             update.Status,
		PaydTransactionRef: update.PaydTransactionRef,
		Message:            update.Message,
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23503" {
				return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "transaction does not exist")
			}
		}

		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to update transaction")
	}

	return &repository.Transaction{
		TransactionID:      transaction.TransactionID,
		PaydTransactionRef: transaction.PaydTransactionRef,
		Message:            transaction.Message,
		UserID:             transaction.UserID,
		Action:             transaction.Action,
		Amount:             transaction.Amount,
		PhoneNumber:        transaction.PhoneNumber,
		NetworkCode:        transaction.NetworkNode,
		Narration:          transaction.Narration,
		Status:             transaction.Status,
		UpdatedAt:          transaction.UpdatedAt,
		CreatedAt:          transaction.CreatedAt,
	}, nil
}
