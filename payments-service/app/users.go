package app

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	db "github.com/EmilioCliff/payment-polling-app/payment-service/db/sqlc"
	"github.com/EmilioCliff/payment-polling-app/payment-service/utils"
	pb "github.com/EmilioCliff/payment-polling-service/shared-grpc/pb"
	"github.com/google/uuid"
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
	TransactionID uuid.UUID `json:"transaction_id"`
	Status        bool      `json:"status"`
	Action        string    `json:"action"`
}

func (app *App) InitiatePayment(ctx context.Context, data initiatePaymentRequest) (initiatePaymentResponse, generalErrorResponse) {
	transactionID, err := uuid.NewRandom()
	if err != nil {
		return initiatePaymentResponse{}, errorHelper("error creating uuid", err, http.StatusInternalServerError, codes.Internal)
	}

	if ctx.Err() == context.Canceled || ctx.Err() == context.DeadlineExceeded {
		return initiatePaymentResponse{}, errorHelper("error creating user", err, http.StatusInternalServerError, codes.Internal)
	}

	transaction, err := app.store.CreateTransaction(ctx, db.CreateTransactionParams{
		TransactionID: transactionID,
		UserID:        data.UserID,
		Action:        data.Action,
		Amount:        int32(data.Amount),
		PhoneNumber:   data.PhoneNumber,
		NetworkNode:   data.NetworkCode,
		Narration:     data.Naration,
	})
	if err != nil {
		return initiatePaymentResponse{}, errorHelper("error creating transaction", err, http.StatusInternalServerError, codes.Internal)
	}

	return initiatePaymentResponse{
		TransactionID: transaction.TransactionID,
		Status:        transaction.Status,
		Action:        transaction.Action,
	}, generalErrorResponse{Status: true}
}

type pollingTransactionRequest struct {
	TransactionId string `json:"transaction_id"`
}

type pollingTransactionResponse struct {
	TransactionID      uuid.UUID `json:"transaction_id"`
	Action             string    `json:"action"`
	Amount             int64     `json:"amount"`
	PhoneNumber        string    `json:"phone_number"`
	NetworkCode        string    `json:"network_code"`
	Naration           string    `json:"naration"`
	Status             bool      `json:"status"`
	PaydUsername       string    `json:"payd_username"`
	PaydUsernameApiKey string    `json:"payd_username_api_key"`
	PaydPasswordApiKey string    `json:"payd_password_api_key"`
}

func (app *App) PollingTransaction(ctx context.Context, transactionID pollingTransactionRequest) (pollingTransactionResponse, generalErrorResponse) {
	log.Println("Log 0:", transactionID)

	transactionUUID, err := uuid.Parse(transactionID.TransactionId)
	if err != nil {
		return pollingTransactionResponse{}, errorHelper("invalid uuid", err, http.StatusInternalServerError, codes.Internal)
	}

	log.Println("Log 1:", transactionUUID)

	transaction, err := app.store.GetTransaction(ctx, transactionUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return pollingTransactionResponse{}, errorHelper("transaction not found", err, http.StatusNotFound, codes.NotFound)
		}
		return pollingTransactionResponse{}, errorHelper("error getting transaction", err, http.StatusInternalServerError, codes.Internal)
	}

	log.Println("Log 2:", transaction)

	userDetails, err := app.authgRPClient.GetUser(ctx, &pb.GetUserRequest{UserId: transaction.UserID})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			grpcCode := st.Code()
			grpcMessage := st.Message()
			return pollingTransactionResponse{}, errorHelper(grpcMessage, err, http.StatusInternalServerError, grpcCode)
		}
	}

	log.Println("Log 3:", userDetails)

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
		Action:             transaction.Action,
		Amount:             int64(transaction.Amount),
		PhoneNumber:        transaction.PhoneNumber,
		NetworkCode:        transaction.NetworkNode,
		Naration:           transaction.Narration,
		Status:             transaction.Status,
		PaydUsername:       userDetails.PaydUsername,
		PaydUsernameApiKey: usernameApiKey,
		PaydPasswordApiKey: passwordApiKey,
	}, generalErrorResponse{Status: true}
}
