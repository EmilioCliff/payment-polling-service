package workers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/EmilioCliff/payment-polling-app/payment-service/internal/repository"
	"github.com/EmilioCliff/payment-polling-app/payment-service/internal/services"
	"github.com/EmilioCliff/payment-polling-app/payment-service/pkg"
	"github.com/hibiken/asynq"
)

const QueueCritical = "critical"

var _ services.TaskProcessor = (*RedisTaskProcessor)(nil)

// type TaskProcessor interface {
// 	Start() error
// 	ProcessPaymentRequestTask(ctx context.Context, task *asynq.Task) error
// 	ProcessWithdrawalRequestTask(ctx context.Context, task *asynq.Task) error
// }

type RedisTaskProcessor struct {
	server *asynq.Server
	config pkg.Config

	paymentUrl    string
	withdrawalUrl string

	TransactionRepository repository.TransactionRepository
}

func NewRedisTaskProcessor(redisOpt *asynq.RedisClientOpt, config pkg.Config) *RedisTaskProcessor {
	server := asynq.NewServer(redisOpt, asynq.Config{
		RetryDelayFunc: asynq.RetryDelayFunc(CustomRetryDelayFunc),
		ErrorHandler:   asynq.ErrorHandlerFunc(ReportError),
		Queues: map[string]int{
			QueueCritical: 10,
		},
	})

	return &RedisTaskProcessor{
		server:        server,
		config:        config,
		paymentUrl:    "https://api.mypayd.app/api/v2/payments",
		withdrawalUrl: "https://api.mypayd.app/api/v2/withdrawal",
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(SendPaymentRequestTask, processor.ProcessPaymentRequestTask)
	mux.HandleFunc(SendWithdrawalRequestTask, processor.ProcessWithdrawalRequestTask)

	return processor.server.Start(mux)
}

func CustomRetryDelayFunc(_ int, _ error, _ *asynq.Task) time.Duration {
	return 500 * time.Millisecond
}

func ReportError(ctx context.Context, task *asynq.Task, err error) {
	retried, _ := asynq.GetRetryCount(ctx)

	maxRetry, _ := asynq.GetMaxRetry(ctx)
	if retried >= maxRetry {
		err = fmt.Errorf("retry exhausted for task %s: %w", task.Type(), err)
	}

	// log it or something
	// errorReportingService.Notify(err)
	log.Println(err)
}
