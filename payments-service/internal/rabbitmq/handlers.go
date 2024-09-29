package rabbitmq

import (
	"context"
	"encoding/json"
	"time"

	"github.com/EmilioCliff/payment-polling-app/payment-service/internal/services"
	"github.com/EmilioCliff/payment-polling-app/payment-service/internal/workers"
	"github.com/EmilioCliff/payment-polling-app/payment-service/pkg"
	"github.com/EmilioCliff/payment-polling-service/shared-grpc/pb"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

type initiatePaymentRequest struct {
	Email       string `json:"email"`
	Action      string `json:"action"`
	Amount      int64  `json:"amount"`
	PhoneNumber string `json:"phone_number"`
	NetworkCode string `json:"network_code"`
	Naration    string `json:"naration"`
}

type initiatePaymentResponse struct {
	TransactionID string `json:"transaction_id"`
	PaymentStatus bool   `json:"payment_status"`
	Action        string `json:"action"`
}

func (r *RabbitConn) handleInitiatePayment(req initiatePaymentRequest) []byte {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	transactionID, err := uuid.NewRandom()
	if err != nil {
		return r.errorRabbitMQResponse(pkg.Errorf(pkg.INTERNAL_ERROR, "failed to create transactionID: %v", err))
	}

	userData, err := r.client.GetUser(ctx, &pb.GetUserRequest{Email: req.Email})
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

	payload := services.SendPaymentWithdrawalRequestPayload{
		TransactionID:      transactionID,
		UserID:             userData.GetUserId(),
		Action:             req.Action,
		Amount:             req.Amount,
		PhoneNumber:        req.PhoneNumber,
		NetworkCode:        req.NetworkCode,
		Naration:           req.Naration,
		PaydUsername:       userData.GetPaydUsername(),
		PaydAccountID:      userData.GetPaydAccountId(),
		PaydPasswordApiKey: passwordApiKey,
		PaydUsernameApiKey: usernameApiKey,
	}

	switch req.Action {
	case "payment":
		err = r.Distributor.DistributeSendPaymentRequestTask(ctx, payload, opts...)
		if err != nil {
			return r.errorRabbitMQResponse(pkg.Errorf(pkg.INTERNAL_ERROR, "failed to distribute payment task: %v", err))
		}

	case "withdrawal":
		err = r.Distributor.DistributeSendWithdrawalRequestTask(ctx, payload, opts...)
		if err != nil {
			return r.errorRabbitMQResponse(pkg.Errorf(pkg.INTERNAL_ERROR, "failed to distribute withdrawal task: %v", err))
		}

	default:
		return r.errorRabbitMQResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid action: %s", req.Action))
	}

	rsp := initiatePaymentResponse{
		TransactionID: transactionID.String(),
		PaymentStatus: false,
		Action:        req.Action,
	}

	rspBytes, err := json.Marshal(rsp)
	if err != nil {
		return r.errorRabbitMQResponse(pkg.Errorf(pkg.INTERNAL_ERROR, "failed to marshal response from initiate-payment %v", err))
	}

	return rspBytes
}

type pollingTransactionRequest struct {
	UserID        int64  `json:"user_id"`
	TransactionId string `json:"transaction_id"`
}

type pollingTransactionResponse struct {
	TransactionID      string `json:"transaction_id"`
	PaydTransactionRef string `json:"payd_transaction_ref"`
	Remarks            string `json:"remarks"`
	Action             string `json:"action"`
	Amount             int32  `json:"amount"`
	PhoneNumber        string `json:"phone_number"`
	NetworkCode        string `json:"network_code"`
	Naration           string `json:"naration"`
	PaymentStatus      bool   `json:"payment_status"`
}

func (r *RabbitConn) handlePollingTransaction(req pollingTransactionRequest) []byte {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	id, err := uuid.Parse(req.TransactionId)
	if err != nil {
		return r.errorRabbitMQResponse(pkg.Errorf(pkg.INVALID_ERROR, "invalid transaction id: %v", err))
	}

	transaction, err := r.TransactionRepository.PollingTransaction(ctx, id)
	if err != nil {
		pkgError, _ := err.(*pkg.Error)

		return r.errorRabbitMQResponse(pkg.Errorf(pkgError.Code, pkgError.Message))
	}

	if transaction.UserID != req.UserID {
		return r.errorRabbitMQResponse(pkg.Errorf(pkg.AUTHENTICATION_ERROR, "cannot access this transaction"))
	}

	rsp := pollingTransactionResponse{
		TransactionID:      transaction.TransactionID.String(),
		PaydTransactionRef: transaction.PaydTransactionRef,
		Remarks:            transaction.Message,
		Action:             transaction.Action,
		Amount:             transaction.Amount,
		PhoneNumber:        transaction.PhoneNumber,
		NetworkCode:        transaction.NetworkCode,
		Naration:           transaction.Narration,
		PaymentStatus:      transaction.Status,
	}

	rspBytes, marshalErr := json.Marshal(rsp)
	if marshalErr != nil {
		return r.errorRabbitMQResponse(pkg.Errorf(pkg.INTERNAL_ERROR, "failed to marshal response from initiate-payment %v", marshalErr))
	}

	return rspBytes
}
