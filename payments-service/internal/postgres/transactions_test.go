package postgres

// package postgres_test

// import (
// 	"context"
// 	"testing"

// 	mockdb "github.com/EmilioCliff/payment-polling-app/payment-service/internal/mock"
// 	"github.com/EmilioCliff/payment-polling-app/payment-service/internal/postgres"
// 	"github.com/EmilioCliff/payment-polling-app/payment-service/internal/rabbitmq"
// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/require"
// 	"go.uber.org/mock/gomock"
// )

// func newTransaction() postgres.InitiatePaymentRequest {
// 	transactionID, _ := uuid.NewRandom()
// 	return postgres.InitiatePaymentRequest{
// 		TransactionID:      transactionID,
// 		UserID:             1,
// 		Action:             "test",
// 		Amount:             100,
// 		PhoneNumber:        "08123456789",
// 		NetworkCode:        "234",
// 		Naration:           "Narration",
// 		PaydTransactionRef: "test",
// 	}
// }

// func TestCreateTransactions(t *testing.T) {
// 	transaction := newTransaction()

// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	store := mockdb.NewMockStore(ctrl)

// 	// r := &rabbitmq.RabbitConn{}

// 	// build stubs
// 	store.EXPECT().
// 		CreateTransactions(gomock.Any(), gomock.Eq(transaction)).
// 		Times(1).
// 		Return(&postgres.InitiatePaymentResponse{}, nil)

// 	transactionResponse, pkgErr := store.CreateTransactions(context.Background(), transaction)
// 	require.Nil(t, pkgErr)
// 	require.Equal(t, transactionResponse.TransactionID, transaction.TransactionID)
// 	require.Equal(t, transactionResponse.Action, transaction.Action)

// 	distributor := mockdb.NewMockTaskDistributor(ctrl)
// }

// func TestGetTransaction(t *testing.T) {}
