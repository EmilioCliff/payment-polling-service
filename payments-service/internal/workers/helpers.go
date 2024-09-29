package workers

import (
	"context"

	"github.com/EmilioCliff/payment-polling-app/payment-service/internal/repository"
	"github.com/hibiken/asynq"
)

func (p *RedisTaskProcessor) createTransaction(
	ctx context.Context,
	req SendPaymentWithdrawalRequestPayload,
	transactionRef string,
	message string,
	status string,
) error {
	if status == "failed" {
		retried, _ := asynq.GetRetryCount(ctx)
		maxRetry, _ := asynq.GetMaxRetry(ctx)

		if retried < maxRetry {
			return nil
		}
	}

	_, err := p.TransactionRepository.CreateTransaction(ctx, repository.Transaction{
		TransactionID:      req.TransactionID,
		PaydTransactionRef: transactionRef,
		Message:            message,
		UserID:             req.UserID,
		Action:             req.Action,
		Amount:             int32(req.Amount),
		PhoneNumber:        req.PhoneNumber,
		NetworkCode:        req.NetworkCode,
		Narration:          req.Naration,
	})
	if err != nil {
		return err
	}

	return nil
}
