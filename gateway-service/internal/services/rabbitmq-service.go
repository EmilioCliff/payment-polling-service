package services

import "github.com/gin-gonic/gin"

type Payload struct {
	Name string `json:"name"`
	Data any    `json:"data"`
}

type RabbitInterface interface {
	RegisterUserViaRabbit(RegisterUserRequest) (int, gin.H)
	LoginUserViaRabbit(LoginUserRequest) (int, gin.H)
	InitiatePaymentViaRabbit(InitiatePaymentRequest) (int, gin.H)
	PollTransactionViaRabbit(PollingTransactionRequest) (int, gin.H)

	SetConsumer(topics []string) error
}
