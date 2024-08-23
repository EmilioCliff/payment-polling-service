package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Payload struct {
	Name string `json:"name"`
	Data any    `json:"data"`
}

type initiatePaymentRequest struct {
	UserID      int64  `json:"user_id" binding:"required"`
	Action      string `json:"action" binding:"required"`
	Amount      int64  `json:"amount" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	NetworkCode string `json:"network_code" binding:"required"`
	Naration    string `json:"naration" binding:"required"`
}

type initiatePaymentResponse struct {
	TransactionID string `json:"transaction_id"`
	PaymentStatus bool   `json:"payment_status"`
	Action        string `json:"action"`
	Message       string `json:"message,omitempty"`
	Status        int    `json:"status,omitempty"`
}

func (server *Server) initiatePaymentViaRabbitMQ(ctx *gin.Context) {
	var req initiatePaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, server.errorResponse(err, "invalid request body"))
		return
	}

	payload := Payload{
		Name: "initiate_payment",
		Data: req,
	}

	payloadRabitData, err := json.Marshal(payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, server.errorResponse(err, "error marshalling request body"))
		return
	}

	correlationID := uuid.New().String()
	responseChannel := make(chan amqp.Delivery, 1)
	defer close(responseChannel)

	server.responseMap.Store(correlationID, responseChannel)
	defer server.responseMap.Delete(correlationID)

	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.amqpChannel.PublishWithContext(c,
		server.config.EXCH,          // exchange
		"payments.initiate_payment", // routing key
		false,                       // mandatory
		false,                       // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: correlationID,
			ReplyTo:       "gateway.initiate_payment",
			Body:          payloadRabitData,
		})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, server.errorResponse(err, "error communicating to the payment service via rabbitmq"))
		return
	}

	select {
	case msg := <-responseChannel:
		if msg.CorrelationId == correlationID {
			var paymentResp initiatePaymentResponse

			err := json.Unmarshal(msg.Body, &paymentResp)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, server.errorResponse(err, "error unmarshalling response"))
				return
			}

			if paymentResp.Message != "" {
				ctx.JSON(paymentResp.Status, gin.H{"error": paymentResp.Message})
				return
			}

			ctx.JSON(http.StatusOK, paymentResp)
		}
	case <-time.After(5 * time.Second):
		ctx.JSON(http.StatusRequestTimeout, gin.H{"error": "Timeout waiting for response"})
	}
}

type pollingTransactionRequest struct {
	TransactionId string `uri:"id" binding:"required"`
}

type pollingTransactionRabbitRequest struct {
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
	Message            string    `json:"message,omitempty"`
	Status             int       `json:"status,omitempty"`
}

func (server *Server) pollTransactionViaRabbitMQ(ctx *gin.Context) {
	var req pollingTransactionRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, server.errorResponse(err, "invalid request body"))
		return
	}

	payload := Payload{
		Name: "polling_transaction",
		Data: pollingTransactionRabbitRequest{
			TransactionId: req.TransactionId,
		},
	}

	payloadRabitData, err := json.Marshal(payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, server.errorResponse(err, "error marshalling request body"))
		return
	}

	correlationID := uuid.New().String()
	responseChannel := make(chan amqp.Delivery, 1)
	defer close(responseChannel)

	server.responseMap.Store(correlationID, responseChannel)
	defer server.responseMap.Delete(correlationID)

	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.amqpChannel.PublishWithContext(c,
		server.config.EXCH,       // exchange
		"payments.poll_payments", // routing key
		false,                    // mandatory
		false,                    // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: correlationID,
			ReplyTo:       "gateway.poll_payments",
			Body:          payloadRabitData,
		})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, server.errorResponse(err, "error communicating to the payment service via rabbitmq"))
		return
	}

	select {
	case msg := <-responseChannel:
		if msg.CorrelationId == correlationID {
			var pollResp pollingTransactionResponse

			err := json.Unmarshal(msg.Body, &pollResp)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, server.errorResponse(err, "error unmarshalling response"))
				return
			}

			if pollResp.Message != "" {
				ctx.JSON(pollResp.Status, gin.H{"error": pollResp.Message})
				return
			}

			ctx.JSON(http.StatusOK, pollResp)
		}
	case <-time.After(5 * time.Second):
		ctx.JSON(http.StatusRequestTimeout, gin.H{"error": "Timeout waiting for response"})
	}
}
