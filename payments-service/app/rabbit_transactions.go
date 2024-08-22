package app

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
)

var (
	ErrorUnmarshalingData = errors.New("error unmarshaling data")
	ErrorMarshalingData   = errors.New("error marshaling data")
)

func (app *App) initiatePaymentHandler(data initiatePaymentRequest) []byte {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	rspData, errRsp := app.InitiatePayment(ctx, data)
	if !errRsp.Status {
		return nil
	}

	rsp, err := json.Marshal(rspData)
	if err != nil {
		return app.errorRabbitMQResponse(ErrorMarshalingData, "error marshaling response", http.StatusInternalServerError)
	}

	return rsp
}

func (app *App) pollingTransactionHandler(data pollingTransactionRequest) []byte {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("polling transaction")

	rspData, errRsp := app.PollingTransaction(ctx, data)
	if !errRsp.Status {
		return nil
	}

	rsp, err := json.Marshal(rspData)
	if err != nil {
		return app.errorRabbitMQResponse(ErrorMarshalingData, "error marshaling response", http.StatusInternalServerError)
	}

	return rsp
}
