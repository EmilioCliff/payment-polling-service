package workers

import (
	"context"

	db "github.com/EmilioCliff/payment-polling-app/payment-service/db/sqlc"
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
	store  *db.Queries
}

func NewRedisTaskProcessor(store *db.Queries, redisOpt *asynq.RedisClientOpt) TaskProcessor {
	server := asynq.NewServer(redisOpt, asynq.Config{
		RetryDelayFunc: asynq.RetryDelayFunc(CustomRetryDelayFunc),
		ErrorHandler:   asynq.ErrorHandlerFunc(ReportError),
		Queues: map[string]int{
			QueueCritical: 10,
		},
	})

	return &RedisTaskProcessor{server: server, store: store}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(SendPaymentRequestTask, processor.ProcessPaymentRequestTask)
	mux.HandleFunc(SendWithdrawalRequestTask, processor.ProcessWithdrawalRequestTask)

	return processor.server.Start(mux)
}
