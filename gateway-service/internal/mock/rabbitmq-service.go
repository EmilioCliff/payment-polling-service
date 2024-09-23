package mock

import (
	"github.com/EmilioCliff/payment-polling-app/gateway-service/internal/services"
	"github.com/gin-gonic/gin"
)

var _ services.RabbitInterface = (*MockRabbitMQService)(nil)

type MockRabbitMQService struct {
	RegisterUserViaRabbitFunc func (services.RegisterUserRequest) (int, gin.H)
	LoginUserViaRabbitFunc func (services.LoginUserRequest) (int, gin.H)
	InitiatePaymentViaRabbitFunc func (services.InitiatePaymentRequest) (int, gin.H)
	PollTransactionViaRabbitFunc func (services.PollingTransactionRequest) (int, gin.H)

	SetConsumerFunc func (topics []string) error
}

func (m *MockRabbitMQService) RegisterUserViaRabbit(req services.RegisterUserRequest) (int, gin.H) {
	return m.RegisterUserViaRabbitFunc(req)
}
func (m *MockRabbitMQService) LoginUserViaRabbit(req services.LoginUserRequest) (int, gin.H) {
	return m.LoginUserViaRabbitFunc(req)
}
func (m *MockRabbitMQService) InitiatePaymentViaRabbit(req services.InitiatePaymentRequest) (int, gin.H) {
	return m.InitiatePaymentViaRabbitFunc(req)
}
func (m *MockRabbitMQService) PollTransactionViaRabbit(req services.PollingTransactionRequest) (int, gin.H) {
	return m.PollTransactionViaRabbitFunc(req)
}
func (m *MockRabbitMQService) SetConsumer(topics []string) error {
	return m.SetConsumerFunc(topics)
}