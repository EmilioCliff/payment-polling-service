package workers

import (
	"context"

	"github.com/EmilioCliff/payment-polling-app/payment-service/internal/postgres"
	"github.com/hibiken/asynq"
)

const QueueCritical = "critical"

type TaskProcessor interface {
	Start() error
	ProcessPaymentRequestTask(ctx context.Context, task *asynq.Task) error
	ProcessWithdrawalRequestTask(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  postgres.Store
}

func NewRedisTaskProcessor(redisOpt *asynq.RedisClientOpt) (TaskProcessor, error) {
	store, err := postgres.GetStore()
	if err != nil {
		return nil, err
	}

	server := asynq.NewServer(redisOpt, asynq.Config{
		RetryDelayFunc: asynq.RetryDelayFunc(CustomRetryDelayFunc),
		ErrorHandler:   asynq.ErrorHandlerFunc(ReportError),
		Queues: map[string]int{
			QueueCritical: 10,
		},
	})

	return &RedisTaskProcessor{server: server, store: store}, nil
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(SendPaymentRequestTask, processor.ProcessPaymentRequestTask)
	mux.HandleFunc(SendWithdrawalRequestTask, processor.ProcessWithdrawalRequestTask)

	return processor.server.Start(mux)
}
