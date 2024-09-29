package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

type SendPaymentWithdrawalRequestPayload struct {
	TransactionID      uuid.UUID `json:"transaction_id"`
	UserID             int64     `json:"user_id"`
	Action             string    `json:"action"`
	Amount             int64     `json:"amount"`
	PhoneNumber        string    `json:"phone_number"`
	NetworkCode        string    `json:"network_code"`
	Naration           string    `json:"naration"`
	PaydUsername       string    `json:"payd_username"`
	PaydAccountID      string    `json:"payd_account_id"`
	PaydPasswordApiKey string    `json:"payd_password_api_key"`
	PaydUsernameApiKey string    `json:"payd_username_api_key"`
}

type TaskProcessor interface {
	Start() error
	ProcessPaymentRequestTask(ctx context.Context, task *asynq.Task) error
	ProcessWithdrawalRequestTask(ctx context.Context, task *asynq.Task) error
}

type TaskDistributor interface {
	DistributeSendPaymentRequestTask(ctx context.Context, payload SendPaymentWithdrawalRequestPayload, opt ...asynq.Option) error
	DistributeSendWithdrawalRequestTask(ctx context.Context, payload SendPaymentWithdrawalRequestPayload, opt ...asynq.Option) error
}
