package mock

import (
	"context"

	"github.com/EmilioCliff/payment-polling-app/payment-service/internal/services"
	"github.com/hibiken/asynq"
)

var _ services.TaskDistributor = (*MockTaskDistributor)(nil)

type MockTaskDistributor struct {
	DistributeSendPaymentRequestTaskFunc    func(ctx context.Context, payload services.SendPaymentWithdrawalRequestPayload, opt ...asynq.Option) error
	DistributeSendWithdrawalRequestTaskFunc func(ctx context.Context, payload services.SendPaymentWithdrawalRequestPayload, opt ...asynq.Option) error
}

func (m *MockTaskDistributor) DistributeSendPaymentRequestTask(
	ctx context.Context,
	payload services.SendPaymentWithdrawalRequestPayload,
	opt ...asynq.Option,
) error {
	return m.DistributeSendPaymentRequestTaskFunc(ctx, payload, opt...)
}

func (m *MockTaskDistributor) DistributeSendWithdrawalRequestTask(
	ctx context.Context,
	payload services.SendPaymentWithdrawalRequestPayload,
	opt ...asynq.Option,
) error {
	return m.DistributeSendWithdrawalRequestTaskFunc(ctx, payload, opt...)
}
