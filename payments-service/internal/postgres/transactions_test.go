package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/EmilioCliff/payment-polling-app/payment-service/internal/postgres/generated"
	mockdb "github.com/EmilioCliff/payment-polling-app/payment-service/internal/postgres/mock"
	"github.com/EmilioCliff/payment-polling-app/payment-service/internal/repository"
	"github.com/EmilioCliff/payment-polling-app/payment-service/pkg"
	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"go.uber.org/mock/gomock"
)

var TestTime = time.Date(2024, time.September, 18, 12, 0, 0, 0, time.UTC)

func NewTestTransactionRepository() *TransactionRepository {
	store := NewStore(pkg.Config{})
	store.conn = nil

	t := NewTransactionService(store)

	return t
}

func TestTransactionRepository_CreateTransaction(t *testing.T) {
	tr := NewTestTransactionRepository()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueries := mockdb.NewMockQuerier(ctrl)

	tr.queries = mockQueries

	tests := []struct {
		name        string
		transaction repository.Transaction
		buildStubs  func(*mockdb.MockQuerier, repository.Transaction)
		wantErr     bool
	}{
		{
			name:        "success",
			transaction: newTransaction(),
			buildStubs: func(q *mockdb.MockQuerier, transaction repository.Transaction) {
				q.EXPECT().CreateTransaction(gomock.Any(), gomock.Eq(generated.CreateTransactionParams{
					TransactionID:      transaction.TransactionID,
					PaydTransactionRef: transaction.PaydTransactionRef,
					Message:            transaction.Message,
					UserID:             transaction.UserID,
					Action:             transaction.Action,
					Amount:             transaction.Amount,
					PhoneNumber:        transaction.PhoneNumber,
					NetworkNode:        transaction.NetworkCode,
					Narration:          transaction.Narration,
				})).Times(1).Return(generated.Transaction{
					TransactionID:      transaction.TransactionID,
					PaydTransactionRef: transaction.PaydTransactionRef,
					Message:            transaction.Message,
					UserID:             transaction.UserID,
					Action:             transaction.Action,
					Amount:             transaction.Amount,
					PhoneNumber:        transaction.PhoneNumber,
					NetworkNode:        transaction.NetworkCode,
					Narration:          transaction.Narration,
					Status:             false,
					UpdatedAt:          TestTime,
					CreatedAt:          TestTime,
				}, nil)
			},
			wantErr: false,
		},
		{
			name:        "missing values",
			transaction: repository.Transaction{},
			buildStubs: func(q *mockdb.MockQuerier, transaction repository.Transaction) {
				q.EXPECT().
					CreateTransaction(gomock.Any(), gomock.Any()).
					Times(0)
			},
			wantErr: true,
		},
		{
			name:        "failed to create transaction",
			transaction: newTransaction(),
			buildStubs: func(q *mockdb.MockQuerier, transaction repository.Transaction) {
				q.EXPECT().
					CreateTransaction(gomock.Any(), gomock.Eq(generated.CreateTransactionParams{
						TransactionID:      transaction.TransactionID,
						PaydTransactionRef: transaction.PaydTransactionRef,
						Message:            transaction.Message,
						UserID:             transaction.UserID,
						Action:             transaction.Action,
						Amount:             transaction.Amount,
						PhoneNumber:        transaction.PhoneNumber,
						NetworkNode:        transaction.NetworkCode,
						Narration:          transaction.Narration,
					})).
					Times(1).
					Return(generated.Transaction{}, errors.New("db internal error"))
			},
			wantErr: true,
		},
		{
			name:        "transaction already exists",
			transaction: newTransaction(),
			buildStubs: func(q *mockdb.MockQuerier, transaction repository.Transaction) {
				q.EXPECT().
					CreateTransaction(gomock.Any(), gomock.Eq(generated.CreateTransactionParams{
						TransactionID:      transaction.TransactionID,
						PaydTransactionRef: transaction.PaydTransactionRef,
						Message:            transaction.Message,
						UserID:             transaction.UserID,
						Action:             transaction.Action,
						Amount:             transaction.Amount,
						PhoneNumber:        transaction.PhoneNumber,
						NetworkNode:        transaction.NetworkCode,
						Narration:          transaction.Narration,
					})).
					Times(1).
					Return(generated.Transaction{}, &pq.Error{
						Code: "23505"})
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(mockQueries, tc.transaction)

			_, err := tr.CreateTransaction(context.Background(), tc.transaction)
			if (err != nil) != tc.wantErr {
				t.Errorf("CreateTransaction() error = %v, wantErr %v", err, tc.wantErr)

				return
			}
		})
	}
}

func TestTransactionRepository_PollingTransaction(t *testing.T) {
	tr := NewTestTransactionRepository()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueries := mockdb.NewMockQuerier(ctrl)

	tr.queries = mockQueries

	tests := []struct {
		name        string
		id          uuid.UUID
		mockQueries func(*mockdb.MockQuerier, uuid.UUID)
		wantErr     bool
	}{
		{
			name: "success",
			id:   uuid.New(),
			mockQueries: func(mockQueries *mockdb.MockQuerier, id uuid.UUID) {
				mockQueries.EXPECT().GetTransaction(gomock.Any(), gomock.Eq(id)).
					Return(generated.Transaction{
						TransactionID:      uuid.New(),
						PaydTransactionRef: gofakeit.UUID(),
						Message:            gofakeit.Sentence(10),
						UserID:             1,
						Action:             "withdrawal",
						Amount:             100,
						PhoneNumber:        gofakeit.Phone(),
						NetworkNode:        "63902",
						Narration:          gofakeit.Sentence(10),
						Status:             true,
					}, nil).Times(1)
			},
			wantErr: false,
		},
		{
			name: "db error",
			id:   uuid.New(),
			mockQueries: func(mockQueries *mockdb.MockQuerier, id uuid.UUID) {
				mockQueries.EXPECT().GetTransaction(gomock.Any(), gomock.Eq(id)).
					Return(generated.Transaction{}, errors.New("db error")).Times(1)
			},
			wantErr: true,
		},
		{
			name: "not found",
			id:   uuid.New(),
			mockQueries: func(mockQueries *mockdb.MockQuerier, id uuid.UUID) {
				mockQueries.EXPECT().GetTransaction(gomock.Any(), gomock.Eq(id)).
					Return(generated.Transaction{}, sql.ErrNoRows).Times(1)
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockQueries(mockQueries, tc.id)

			_, err := tr.PollingTransaction(context.Background(), tc.id)
			if (err != nil) != tc.wantErr {
				t.Errorf("PollingTransaction() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestTransactionRepository_UpdateTransaction(t *testing.T) {
	tr := NewTestTransactionRepository()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueries := mockdb.NewMockQuerier(ctrl)

	tr.queries = mockQueries

	tests := []struct {
		name        string
		id          uuid.UUID
		transaction repository.TransactionUpdate
		buildStubs  func(*mockdb.MockQuerier, uuid.UUID, repository.TransactionUpdate)
		wantErr     bool
	}{
		{
			name: "success",
			id:   uuid.New(),
			transaction: repository.TransactionUpdate{
				PaydTransactionRef: gofakeit.UUID(),
				Message:            gofakeit.Sentence(10),
				Status:             true,
			},
			buildStubs: func(q *mockdb.MockQuerier, id uuid.UUID, transaction repository.TransactionUpdate) {
				q.EXPECT().
					UpdateTransaction(gomock.Any(), gomock.Eq(generated.UpdateTransactionParams{
						TransactionID:      id,
						PaydTransactionRef: transaction.PaydTransactionRef,
						Message:            transaction.Message,
						Status:             transaction.Status,
					})).
					Times(1).
					Return(generated.Transaction{
						TransactionID:      uuid.New(),
						PaydTransactionRef: transaction.PaydTransactionRef,
						UserID:             1,
						Message:            transaction.Message,
						Action:             "withdrawal",
						Amount:             100,
						PhoneNumber:        gofakeit.Phone(),
						NetworkNode:        "63902",
						Narration:          gofakeit.Sentence(10),
						Status:             transaction.Status,
						UpdatedAt:          TestTime,
						CreatedAt:          TestTime,
					}, nil)
			},
			wantErr: false,
		},
		{
			name: "db error",
			id:   uuid.New(),
			transaction: repository.TransactionUpdate{
				PaydTransactionRef: gofakeit.UUID(),
				Message:            gofakeit.Sentence(10),
				Status:             true,
			},
			buildStubs: func(q *mockdb.MockQuerier, id uuid.UUID, transaction repository.TransactionUpdate) {
				q.EXPECT().
					UpdateTransaction(gomock.Any(), gomock.Eq(generated.UpdateTransactionParams{
						TransactionID:      id,
						PaydTransactionRef: transaction.PaydTransactionRef,
						Message:            transaction.Message,
						Status:             transaction.Status,
					})).
					Return(generated.Transaction{}, errors.New("db error")).
					Times(1)
			},
			wantErr: true,
		},
		{
			name: "transaction not found",
			id:   uuid.New(),
			transaction: repository.TransactionUpdate{
				PaydTransactionRef: gofakeit.UUID(),
				Message:            gofakeit.Sentence(10),
				Status:             true,
			},
			buildStubs: func(q *mockdb.MockQuerier, id uuid.UUID, transaction repository.TransactionUpdate) {
				q.EXPECT().
					UpdateTransaction(gomock.Any(), gomock.Eq(generated.UpdateTransactionParams{
						TransactionID:      id,
						PaydTransactionRef: transaction.PaydTransactionRef,
						Message:            transaction.Message,
						Status:             transaction.Status,
					})).
					Times(1).
					Return(generated.Transaction{}, &pq.Error{Code: "23503"})
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(mockQueries, tc.id, tc.transaction)

			_, err := tr.UpdateTransaction(context.Background(), tc.id, tc.transaction)
			if (err != nil) != tc.wantErr {
				t.Errorf("UpdateTransaction() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func newTransaction() repository.Transaction {
	return repository.Transaction{
		TransactionID:      uuid.New(),
		PaydTransactionRef: gofakeit.UUID(),
		UserID:             1,
		Message:            gofakeit.Sentence(10),
		Action:             "withdrawal",
		Amount:             100,
		PhoneNumber:        gofakeit.Phone(),
		NetworkCode:        "63902",
		Narration:          gofakeit.Sentence(10),
		Status:             false,
	}
}
