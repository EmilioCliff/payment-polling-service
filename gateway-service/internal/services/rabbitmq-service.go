package services

type Payload struct {
	Name string `json:"name"`
	Data any    `json:"data"`
}

type RabbitInterface interface {
	RegisterUserViaRabbit(RegisterUserRequest) (int, RegisterUserResponse)
	LoginUserViaRabbit(LoginUserRequest) (int, LoginUserResponse)
	InitiatePaymentViaRabbit(InitiatePaymentRequest) (int, InitiatePaymentResponse)
	PollTransactionViaRabbit(PollingTransactionRequest) (int, PollingTransactionResponse)

	SetConsumer([]string, chan struct{}) error
}
