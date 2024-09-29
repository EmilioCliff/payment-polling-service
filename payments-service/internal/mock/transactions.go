package mock

import (
	"context"

	"github.com/EmilioCliff/payment-polling-app/payment-service/internal/repository"
	"github.com/google/uuid"
)

var _ repository.TransactionRepository = (*MockTransactionRepository)(nil)

type MockTransactionRepository struct {
	CreateTransactionFunc  func(context.Context, repository.Transaction) (*repository.Transaction, error)
	PollingTransactionFunc func(context.Context, uuid.UUID) (*repository.Transaction, error)
	UpdateTransactionFunc  func(context.Context, uuid.UUID, repository.TransactionUpdate) (*repository.Transaction, error)
}

func (m *MockTransactionRepository) CreateTransaction(ctx context.Context, req repository.Transaction) (*repository.Transaction, error) {
	return m.CreateTransactionFunc(ctx, req)
}

func (m *MockTransactionRepository) PollingTransaction(ctx context.Context, id uuid.UUID) (*repository.Transaction, error) {
	return m.PollingTransactionFunc(ctx, id)
}

func (m *MockTransactionRepository) UpdateTransaction(
	ctx context.Context,
	id uuid.UUID,
	req repository.TransactionUpdate,
) (*repository.Transaction, error) {
	return m.UpdateTransactionFunc(ctx, id, req)
}
