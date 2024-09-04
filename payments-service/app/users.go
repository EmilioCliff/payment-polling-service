package app

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/EmilioCliff/payment-polling-app/payment-service/utils"
	"github.com/EmilioCliff/payment-polling-app/payment-service/workers"
	pb "github.com/EmilioCliff/payment-polling-service/shared-grpc/pb"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type initiatePaymentRequest struct {
	UserID      int64  `json:"user_id"`
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

func (app *App) InitiatePayment(ctx context.Context, data initiatePaymentRequest) (initiatePaymentResponse, generalErrorResponse) {
	userDetails, err := app.authgRPClient.GetUser(ctx, &pb.GetUserRequest{UserId: data.UserID})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			grpcCode := st.Code()
			grpcMessage := st.Message()
			return initiatePaymentResponse{}, errorHelper(grpcMessage, err, http.StatusInternalServerError, grpcCode)
		}
	}

	passwordApiKey, err := utils.Decrypt(userDetails.PaydPasswordKey, []byte(app.config.ENCRYPTION_KEY))
	if err != nil {
		return initiatePaymentResponse{}, errorHelper("error decrypting password api key", err, http.StatusInternalServerError, codes.Internal)
	}

	usernameApiKey, err := utils.Decrypt(userDetails.PaydUsernameKey, []byte(app.config.ENCRYPTION_KEY))
	if err != nil {
		return initiatePaymentResponse{}, errorHelper("error decrypting username api key", err, http.StatusInternalServerError, codes.Internal)
	}

	log.Println(usernameApiKey)
	log.Println(passwordApiKey)

	opts := []asynq.Option{
		asynq.MaxRetry(1),
		asynq.Queue(workers.QueueCritical),
	}

	transactionID, err := uuid.NewRandom()
	if err != nil {
		return initiatePaymentResponse{}, errorHelper("error creating uuid", err, http.StatusInternalServerError, codes.Internal)
	}

	err = app.distributor.DistributeSendPaymentRequestTask(ctx, workers.SendPaymentWithdrawalRequestPayload{
		TransactionID:      transactionID,
		UserID:             data.UserID,
		Action:             data.Action,
		Amount:             data.Amount,
		PhoneNumber:        data.PhoneNumber,
		NetworkCode:        data.NetworkCode,
		Naration:           data.Naration,
		PaydUsername:       userDetails.GetPaydUsername(),
		PaydPasswordApiKey: passwordApiKey,
		PaydUsernameApiKey: usernameApiKey,
	}, opts...)
	if err != nil {
		return initiatePaymentResponse{}, errorHelper("error distributing task", err, http.StatusInternalServerError, codes.Internal)
	}

	return initiatePaymentResponse{
		TransactionID: transactionID.String(),
		PaymentStatus: false,
		Action:        data.Action,
	}, generalErrorResponse{Status: true}
}

type pollingTransactionRequest struct {
	TransactionId string `json:"transaction_id"`
}

type pollingTransactionResponse struct {
	TransactionID      uuid.UUID `json:"transaction_id"`
	PaydTransactionRef string    `json:"payd_transaction_ref"`
	Action             string    `json:"action"`
	Amount             int64     `json:"amount"`
	PhoneNumber        string    `json:"phone_number"`
	NetworkCode        string    `json:"network_code"`
	Naration           string    `json:"naration"`
	PaymentStatus      bool      `json:"payment_status"`
	PaydUsername       string    `json:"payd_username"`
	PaydUsernameApiKey string    `json:"payd_username_api_key"`
	PaydPasswordApiKey string    `json:"payd_password_api_key"`
}

func (app *App) PollingTransaction(ctx context.Context, transactionID pollingTransactionRequest) (pollingTransactionResponse, generalErrorResponse) {
	transactionUUID, err := uuid.Parse(transactionID.TransactionId)
	if err != nil {
		return pollingTransactionResponse{}, errorHelper("invalid uuid", err, http.StatusInternalServerError, codes.Internal)
	}

	transaction, err := app.store.GetTransaction(ctx, transactionUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return pollingTransactionResponse{}, errorHelper("transaction not found", err, http.StatusNotFound, codes.NotFound)
		}
		return pollingTransactionResponse{}, errorHelper("error getting transaction", err, http.StatusInternalServerError, codes.Internal)
	}

	userDetails, err := app.authgRPClient.GetUser(ctx, &pb.GetUserRequest{UserId: transaction.UserID})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			grpcCode := st.Code()
			grpcMessage := st.Message()
			return pollingTransactionResponse{}, errorHelper(grpcMessage, err, http.StatusInternalServerError, grpcCode)
		}
	}

	passwordApiKey, err := utils.Decrypt(userDetails.PaydPasswordKey, []byte(app.config.ENCRYPTION_KEY))
	if err != nil {
		return pollingTransactionResponse{}, errorHelper("error decrypting password api key", err, http.StatusInternalServerError, codes.Internal)
	}

	usernameApiKey, err := utils.Decrypt(userDetails.PaydUsernameKey, []byte(app.config.ENCRYPTION_KEY))
	if err != nil {
		return pollingTransactionResponse{}, errorHelper("error decrypting username api key", err, http.StatusInternalServerError, codes.Internal)
	}

	return pollingTransactionResponse{
		TransactionID:      transaction.TransactionID,
		PaydTransactionRef: transaction.PaydTransactionRef,
		Action:             transaction.Action,
		Amount:             int64(transaction.Amount),
		PhoneNumber:        transaction.PhoneNumber,
		NetworkCode:        transaction.NetworkNode,
		Naration:           transaction.Narration,
		PaymentStatus:      transaction.Status,
		PaydUsername:       userDetails.PaydUsername,
		PaydUsernameApiKey: usernameApiKey,
		PaydPasswordApiKey: passwordApiKey,
	}, generalErrorResponse{Status: true}
}

type callBackUriRequest struct {
	ID string `uri:"id" binding:"required"`
}

func (app *App) callBack(ctx *gin.Context) {
	var uri callBackUriRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error_message": "bad_request - missing transaction id"})
		return
	}

	_, err := uuid.Parse(uri.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error_message": "invalid url"})
		return
	}

	var req any
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error_message": "no body received"})
		return
	}

	// probably change the tranasction status after receiving the callback

	log.Println("request: ", req)
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}
