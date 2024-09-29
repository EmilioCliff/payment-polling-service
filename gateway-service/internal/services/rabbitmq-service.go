package services

type Payload struct {
	Name string `json:"name"`
	Data []byte `json:"data"`
	// Data any    `json:"data"`
}

type RabbitInterface interface {
	RegisterUserViaRabbit(RegisterUserRequest) (int, RegisterUserResponse)
	LoginUserViaRabbit(LoginUserRequest) (int, LoginUserResponse)
	InitiatePaymentViaRabbit(InitiatePaymentRequest) (int, InitiatePaymentResponse)
	PollTransactionViaRabbit(PollingTransactionRequest, int64) (int, PollingTransactionResponse)

	SetConsumer([]string, chan struct{}) error
}
