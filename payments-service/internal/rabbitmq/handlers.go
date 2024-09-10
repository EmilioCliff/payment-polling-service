package rabbitmq

import (
	"context"
	"encoding/json"
	"time"

	"github.com/EmilioCliff/payment-polling-app/payment-service/internal/postgres"
	"github.com/EmilioCliff/payment-polling-app/payment-service/internal/workers"
	"github.com/EmilioCliff/payment-polling-app/payment-service/pkg"
	"github.com/EmilioCliff/payment-polling-service/shared-grpc/pb"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

type initiatePaymentRequest struct {
	UserID      int64  `json:"user_id" binding:"required"`
	Action      string `json:"action" binding:"required"`
	Amount      int64  `json:"amount" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	NetworkCode string `json:"network_code" binding:"required"`
	Naration    string `json:"naration" binding:"required"`
}

type InitiatePaymentResponse struct {
	TransactionID string `json:"transaction_id"`
	PaymentStatus bool   `json:"payment_status"`
	Action        string `json:"action"`
	Message       string `json:"message,omitempty"`
	Status        int    `json:"status,omitempty"`
}

func (r *RabbitConn) handleInitiatePayment(req initiatePaymentRequest) []byte {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	transactionID, err := uuid.NewRandom()
	if err != nil {
		return r.errorRabbitMQResponse(pkg.Errorf(pkg.INTERNAL_ERROR, "failed to create transactionID: %v", err))
	}

	userData, err := r.client.GetUser(ctx, &pb.GetUserRequest{UserId: req.UserID})
	if err != nil {
		return r.errorRabbitMQResponse(pkg.Errorf(pkg.INTERNAL_ERROR, "failed to get user data from auth: %v", err))
	}

	opts := []asynq.Option{
		asynq.MaxRetry(1),
		asynq.Queue(workers.QueueCritical),
	}

	passwordApiKey, err := pkg.Decrypt(userData.GetPaydPasswordKey(), []byte(r.config.ENCRYPTION_KEY))
	if err != nil {
		return r.errorRabbitMQResponse(pkg.Errorf(pkg.INTERNAL_ERROR, "failed to decrypt payd password key: %v", err))
	}

	usernameApiKey, err := pkg.Decrypt(userData.GetPaydUsernameKey(), []byte(r.config.ENCRYPTION_KEY))
	if err != nil {
		return r.errorRabbitMQResponse(pkg.Errorf(pkg.INTERNAL_ERROR, "failed to decrypt payd username key: %v", err))
	}

	payload := workers.SendPaymentWithdrawalRequestPayload{
		TransactionID:      transactionID,
		UserID:             req.UserID,
		Action:             req.Action,
		Amount:             req.Amount,
		PhoneNumber:        req.PhoneNumber,
		NetworkCode:        req.NetworkCode,
		Naration:           req.Naration,
		PaydUsername:       userData.GetPaydUsername(),
		PaydPasswordApiKey: passwordApiKey,
		PaydUsernameApiKey: usernameApiKey,
	}

	switch req.Action {
	case "payment":
		err = r.distributor.DistributeSendPaymentRequestTask(ctx, payload, opts...)
		if err != nil {
			return r.errorRabbitMQResponse(pkg.Errorf(pkg.INTERNAL_ERROR, "failed to distribute payment task: %v", err))
		}

	case "withdrawal":
		err = r.distributor.DistributeSendWithdrawalRequestTask(ctx, payload, opts...)
		if err != nil {
			return r.errorRabbitMQResponse(pkg.Errorf(pkg.INTERNAL_ERROR, "failed to distribute payment task: %v", err))
		}

	default:
		return r.errorRabbitMQResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid action: %s", req.Action))
	}

	rsp := postgres.InitiatePaymentResponse{
		TransactionID: transactionID,
		PaymentStatus: false,
		Action:        req.Action,
		Amount:        int32(req.Amount),
	}

	rspBytes, err := json.Marshal(rsp)
	if err != nil {
		return r.errorRabbitMQResponse(pkg.Errorf(pkg.INTERNAL_ERROR, "failed to marshal response from initiate-payment %v", err))
	}

	return rspBytes
}

func (r *RabbitConn) handlePollingTransaction(req postgres.PollingTransactionRequest) []byte {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	rsp, err := r.store.PollingTransaction(ctx, req)
	if err != nil {
		return r.errorRabbitMQResponse(pkg.Errorf(err.Code, "failed to get transaction in rebbit: %v", err.Message))
	}

	rspBytes, marshalErr := json.Marshal(rsp)
	if marshalErr != nil {
		return r.errorRabbitMQResponse(pkg.Errorf(pkg.INTERNAL_ERROR, "failed to marshal response from initiate-payment %v", marshalErr))
	}

	return rspBytes
}
