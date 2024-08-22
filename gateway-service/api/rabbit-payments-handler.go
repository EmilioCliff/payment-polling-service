package api

import (
	"context"
	"encoding/json"
	"log"
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
			ctx.JSON(http.StatusOK, gin.H{"payment response": string(msg.Body)})
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

	log.Printf("%v", payload)

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
			ctx.JSON(http.StatusOK, gin.H{"payment response": string(msg.Body)})
		}
	case <-time.After(5 * time.Second):
		ctx.JSON(http.StatusRequestTimeout, gin.H{"error": "Timeout waiting for response"})
	}
}
