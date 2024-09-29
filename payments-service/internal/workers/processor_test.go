package workers

import (
	"github.com/EmilioCliff/payment-polling-app/payment-service/internal/mock"
	"github.com/EmilioCliff/payment-polling-app/payment-service/pkg"
	"github.com/hibiken/asynq"
)

type TestRedisProcessor struct {
	redisProcessor *RedisTaskProcessor

	TransactionRepository mock.MockTransactionRepository
}

func NewTestRedisProcessor() *TestRedisProcessor {
	p := &TestRedisProcessor{
		redisProcessor: NewRedisTaskProcessor(&asynq.RedisClientOpt{}, pkg.Config{}),
	}

	p.redisProcessor.TransactionRepository = &p.TransactionRepository

	return p
}
