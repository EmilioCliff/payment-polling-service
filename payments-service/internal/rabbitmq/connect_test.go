package rabbitmq

import (
	"github.com/EmilioCliff/payment-polling-app/payment-service/internal/mock"
	"github.com/EmilioCliff/payment-polling-app/payment-service/pkg"
	// "github.com/EmilioCliff/payment-polling-service/shared-grpc/mockpb"
	// "go.uber.org/mock/gomock"
)

type TestRabbitHandler struct {
	rabbit *RabbitConn

	TransactionRepository mock.MockTransactionRepository
	TastDistributor       mock.MockTaskDistributor
}

func NewTestRabbitHandler() *TestRabbitHandler {
	rt := &TestRabbitHandler{
		rabbit: NewRabbitConn(pkg.Config{
			ENCRYPTION_KEY: "12345678901234567890123456789012",
		}, nil),
	}

	rt.rabbit.TransactionRepository = &rt.TransactionRepository
	rt.rabbit.Distributor = &rt.TastDistributor

	return rt
}
