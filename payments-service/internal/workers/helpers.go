package workers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

type createTransactionData struct {
	TransactionID      uuid.UUID `json:"transaction_id"`
	PaydTransactionRef string    `json:"payd_transaction_ref"`
	UserID             int64     `json:"user_id"`
	Action             string    `json:"action"`
	Amount             int64     `json:"amount"`
	PhoneNumber        string    `json:"phone_number"`
	NetworkCode        string    `json:"network_code"`
	Naration           string    `json:"naration"`
}

// RetryDelayFunc calculates the retry delay duration for a failed task given
// the retry count, error, and the task.
//
// n is the number of times the task has been retried.
// e is the error returned by the task handler.
// t is the task in question.
// type RetryDelayFunc func(n int, e error, t *Task) time.Duration

func CustomRetryDelayFunc(_ int, _ error, _ *asynq.Task) time.Duration {
	return 500 * time.Millisecond
}

func ReportError(ctx context.Context, task *asynq.Task, err error) {
	retried, _ := asynq.GetRetryCount(ctx)
	maxRetry, _ := asynq.GetMaxRetry(ctx)
	if retried >= maxRetry {
		err = fmt.Errorf("retry exhausted for task %s: %w", task.Type(), err)
	}
	log.Println(err)
	// log it or something
	// errorReportingService.Notify(err)
}
