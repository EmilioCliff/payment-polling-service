package rabbitmq

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/EmilioCliff/payment-polling-app/gateway-service/internal/services"
	"github.com/EmilioCliff/payment-polling-app/gateway-service/pkg"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (r *RabbitHandler) InitiatePaymentViaRabbit(req services.InitiatePaymentRequest) (int, gin.H) {
	payload := services.Payload{
		Name: "initiate_payment",
		Data: req,
	}

	payloadRabitData, err := json.Marshal(payload)
	if err != nil {
		return http.StatusInternalServerError, pkg.ErrorResponse(err, "error marshalling request body")
	}

	correlationID := uuid.New().String()

	responseChannel := make(chan amqp.Delivery, 1)
	defer close(responseChannel)

	r.RspMap.Set(correlationID, responseChannel)
	defer r.RspMap.Delete(correlationID)

	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = r.Channel.PublishWithContext(c,
		r.config.EXCH,               // exchange
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
		return http.StatusInternalServerError, pkg.ErrorResponse(
			err,
			"error communicating to the payment service via rabbitmq",
		)
	}

	select {
	case msg := <-responseChannel:
		if msg.CorrelationId == correlationID {
			var paymentResp services.InitiatePaymentResponse

			err := json.Unmarshal(msg.Body, &paymentResp)
			if err != nil {
				return http.StatusInternalServerError, pkg.ErrorResponse(err, "error unmarshalling response")
			}

			if paymentResp.Message != "" {
				return paymentResp.StatusCode, gin.H{"error": paymentResp.Message}
			}

			return http.StatusOK, gin.H{"response": paymentResp}
		}
	case <-time.After(5 * time.Second):
		return http.StatusRequestTimeout, gin.H{"error": "timeout waiting for response"}
	}

	return http.StatusInternalServerError, gin.H{"error": "unknown error"}
}

type pollingTransactionRabbitRequest struct {
	TransactionId string `json:"transaction_id"`
}

func (r *RabbitHandler) PollTransactionViaRabbit(req services.PollingTransactionRequest) (int, gin.H) {
	payload := services.Payload{
		Name: "polling_transaction",
		Data: pollingTransactionRabbitRequest{
			TransactionId: req.TransactionId,
		},
	}

	payloadRabitData, err := json.Marshal(payload)
	if err != nil {
		return http.StatusInternalServerError, pkg.ErrorResponse(err, "error marshalling request body")
	}

	correlationID := uuid.New().String()

	responseChannel := make(chan amqp.Delivery, 1)
	defer close(responseChannel)

	r.RspMap.Set(correlationID, responseChannel)
	defer r.RspMap.Delete(correlationID)

	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = r.Channel.PublishWithContext(c,
		r.config.EXCH,            // exchange
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
		return http.StatusInternalServerError, pkg.ErrorResponse(
			err,
			"error communicating to the payment service via rabbitmq",
		)
	}

	select {
	case msg := <-responseChannel:
		if msg.CorrelationId == correlationID {
			var pollResp services.PollingTransactionResponse

			err := json.Unmarshal(msg.Body, &pollResp)
			if err != nil {
				return http.StatusInternalServerError, pkg.ErrorResponse(err, "error unmarshalling response")
			}

			if pollResp.Message != "" {
				return pollResp.StatusCode, gin.H{"error": pollResp.Message}
			}

			return http.StatusOK, gin.H{"response": pollResp}
		}
	case <-time.After(5 * time.Second):
		return http.StatusRequestTimeout, gin.H{"error": "timeout waiting for response"}
	}

	return http.StatusInternalServerError, gin.H{"error": "unknown error"}
}
