package mock

import (
	"github.com/EmilioCliff/payment-polling-app/gateway-service/internal/services"
)

var _ services.RabbitInterface = (*MockRabbitMQService)(nil)

type MockRabbitMQService struct {
	RegisterUserViaRabbitFunc    func(services.RegisterUserRequest) (int, services.RegisterUserResponse)
	LoginUserViaRabbitFunc       func(services.LoginUserRequest) (int, services.LoginUserResponse)
	InitiatePaymentViaRabbitFunc func(services.InitiatePaymentRequest) (int, services.InitiatePaymentResponse)
	PollTransactionViaRabbitFunc func(services.PollingTransactionRequest, int64) (int, services.PollingTransactionResponse)

	SetConsumerFunc func(topics []string) error
}

func (m *MockRabbitMQService) RegisterUserViaRabbit(req services.RegisterUserRequest) (int, services.RegisterUserResponse) {
	return m.RegisterUserViaRabbitFunc(req)
}

func (m *MockRabbitMQService) LoginUserViaRabbit(req services.LoginUserRequest) (int, services.LoginUserResponse) {
	return m.LoginUserViaRabbitFunc(req)
}

func (m *MockRabbitMQService) InitiatePaymentViaRabbit(req services.InitiatePaymentRequest) (int, services.InitiatePaymentResponse) {
	return m.InitiatePaymentViaRabbitFunc(req)
}

func (m *MockRabbitMQService) PollTransactionViaRabbit(
	req services.PollingTransactionRequest,
	userID int64,
) (int, services.PollingTransactionResponse) {
	return m.PollTransactionViaRabbitFunc(req, userID)
}

func (m *MockRabbitMQService) SetConsumer(topics []string, _ chan struct{}) error {
	return m.SetConsumerFunc(topics)
}
