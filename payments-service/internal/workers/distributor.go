package workers

import (
	"github.com/EmilioCliff/payment-polling-app/payment-service/internal/services"
	"github.com/hibiken/asynq"
)

// type TaskDistributor interface {
// 	DistributeSendPaymentRequestTask(ctx context.Context, payload SendPaymentWithdrawalRequestPayload, opt ...asynq.Option) error
// 	DistributeSendWithdrawalRequestTask(ctx context.Context, payload SendPaymentWithdrawalRequestPayload, opt ...asynq.Option) error
// }

var _ services.TaskDistributor = (*RedisTaskDistributor)(nil)

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisOpt *asynq.RedisClientOpt) *RedisTaskDistributor {
	client := asynq.NewClient(redisOpt)

	return &RedisTaskDistributor{
		client: client,
	}
}
