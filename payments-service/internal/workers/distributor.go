package workers

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	DistributeSendPaymentRequestTask(ctx context.Context, payload SendPaymentWithdrawalRequestPayload, opt ...asynq.Option) error
	DistributeSendWithdrawalRequestTask(ctx context.Context, payload SendPaymentWithdrawalRequestPayload, opt ...asynq.Option) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisOpt *asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOpt)

	return &RedisTaskDistributor{
		client: client,
	}
}
