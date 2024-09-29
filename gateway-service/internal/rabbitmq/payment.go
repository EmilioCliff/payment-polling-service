package rabbitmq

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/EmilioCliff/payment-polling-app/gateway-service/internal/services"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (r *RabbitHandler) InitiatePaymentViaRabbit(req services.InitiatePaymentRequest) (int, services.InitiatePaymentResponse) {
	dataBytes, err := json.Marshal(req)
	if err != nil {
		return http.StatusInternalServerError, services.InitiatePaymentResponse{Message: "internal error", StatusCode: http.StatusInternalServerError}
	}

	payload := services.Payload{
		Name: "initiate_payment",
		Data: dataBytes,
	}

	payloadRabitData, err := json.Marshal(payload)
	if err != nil {
		return http.StatusInternalServerError, services.InitiatePaymentResponse{Message: "internal error", StatusCode: http.StatusInternalServerError}
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
		return http.StatusInternalServerError, services.InitiatePaymentResponse{Message: "internal error", StatusCode: http.StatusInternalServerError}
	}

	select {
	case msg := <-responseChannel:
		if msg.CorrelationId == correlationID {
			var paymentResp services.InitiatePaymentResponse

			err := json.Unmarshal(msg.Body, &paymentResp)
			if err != nil {
				return http.StatusInternalServerError, services.InitiatePaymentResponse{
					Message:    "internal error",
					StatusCode: http.StatusInternalServerError,
				}
			}

			if paymentResp.Message != "" {
				return paymentResp.StatusCode, services.InitiatePaymentResponse{Message: paymentResp.Message, StatusCode: paymentResp.StatusCode}
			}

			return http.StatusOK, paymentResp
		}
	case <-time.After(5 * time.Second):
		return http.StatusRequestTimeout, services.InitiatePaymentResponse{
			Message:    "timeout waiting for response. Try again",
			StatusCode: http.StatusInternalServerError,
		}
	}

	return http.StatusInternalServerError, services.InitiatePaymentResponse{Message: "internal error", StatusCode: http.StatusInternalServerError}
}

type pollingTransactionRabbitRequest struct {
	UserID        int64  `json:"user_id"`
	TransactionId string `json:"transaction_id"`
}

func (r *RabbitHandler) PollTransactionViaRabbit(req services.PollingTransactionRequest, userID int64) (int, services.PollingTransactionResponse) {
	dataBytes, err := json.Marshal(pollingTransactionRabbitRequest{
		UserID:        userID,
		TransactionId: req.TransactionId,
	})
	if err != nil {
		return http.StatusInternalServerError, services.PollingTransactionResponse{
			Message:    "internal error",
			StatusCode: http.StatusInternalServerError,
		}
	}

	payload := services.Payload{
		Name: "polling_transaction",
		Data: dataBytes,
	}

	payloadRabitData, err := json.Marshal(payload)
	if err != nil {
		return http.StatusInternalServerError, services.PollingTransactionResponse{
			Message:    "internal error",
			StatusCode: http.StatusInternalServerError,
		}
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
		return http.StatusInternalServerError, services.PollingTransactionResponse{
			Message:    "internal error",
			StatusCode: http.StatusInternalServerError,
		}
	}

	select {
	case msg := <-responseChannel:
		if msg.CorrelationId == correlationID {
			var pollResp services.PollingTransactionResponse

			err := json.Unmarshal(msg.Body, &pollResp)
			if err != nil {
				return http.StatusInternalServerError, services.PollingTransactionResponse{
					Message:    "internal error",
					StatusCode: http.StatusInternalServerError,
				}
			}

			if pollResp.Message != "" {
				return pollResp.StatusCode, services.PollingTransactionResponse{Message: pollResp.Message, StatusCode: pollResp.StatusCode}
			}

			return http.StatusOK, pollResp
		}
	case <-time.After(5 * time.Second):
		return http.StatusRequestTimeout, services.PollingTransactionResponse{
			Message:    "timeout waiting for response. Try again",
			StatusCode: http.StatusInternalServerError,
		}
	}

	return http.StatusInternalServerError, services.PollingTransactionResponse{Message: "internal error", StatusCode: http.StatusInternalServerError}
}
