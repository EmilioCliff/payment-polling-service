package rabbitmq

import (
	"encoding/json"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/EmilioCliff/payment-polling-app/gateway-service/internal/services"
	"github.com/brianvoe/gofakeit"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/require"
)

func TestRabbitHandler_InitiatePaymentViaRabbit(t *testing.T) {
	testRabbit, err := NewTestRabbitHandler()
	require.NoError(t, err)

	defer func() {
		// close the channel and terminate the container
		testRabbit.rabbit.Channel.Close()

		if err := testRabbit.container.Terminate(testRabbit.ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()

	req := services.InitiatePaymentRequest{
		UserID:      1,
		Action:      "withdrawal",
		Amount:      100,
		PhoneNumber: gofakeit.Phone(),
		NetworkCode: gofakeit.Word(),
		Naration:    gofakeit.Name(),
	}

	rsp := services.InitiatePaymentResponse{
		TransactionID: gofakeit.UUID(),
		PaymentStatus: true,
		Action:        req.Action,
	}

	rspBytes, err := json.Marshal(rsp)
	require.NoError(t, err)

	tests := []struct {
		name           string
		mockDelivery   amqp.Delivery
		expectedStatus int
		expectedMsg    gin.H
		sleep          time.Duration
	}{
		{
			name: "success",
			mockDelivery: amqp.Delivery{
				ContentType: "text/plain",
				Body:        rspBytes,
			},
			expectedStatus: http.StatusOK,
			expectedMsg:    gin.H{"response": rsp},
			sleep:          100 * time.Millisecond,
		},
		{
			name: "time out",
			mockDelivery: amqp.Delivery{
				ContentType: "text/plain",
				Body:        rspBytes,
			},
			expectedStatus: http.StatusRequestTimeout,
			expectedMsg:    gin.H{"error": "timeout waiting for response"},
			sleep:          6 * time.Second,
		},
	}

	for index, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			statusCodeChan := make(chan int, 1)
			defer close(statusCodeChan)

			msgChan := make(chan gin.H, 1)
			defer close(msgChan)

			go func(statusCodeChan chan int, msgChan chan gin.H) {
				statusCode, msg := testRabbit.rabbit.InitiatePaymentViaRabbit(req)
				statusCodeChan <- statusCode
				msgChan <- msg
			}(statusCodeChan, msgChan)

			// sleep so that the goroutine can send the message
			time.Sleep(tc.sleep)

			var i int
			for correlationID, responseChannel := range testRabbit.rabbit.RspMap.data {
				if i == index {
					log.Printf("CorrelationID for test case '%s': %s", tc.name, correlationID)
					tc.mockDelivery.CorrelationId = correlationID
					responseChannel <- tc.mockDelivery
				}

				i++
			}

			code := <-statusCodeChan
			msg := <-msgChan

			require.Equal(t, tc.expectedStatus, code)
			require.Equal(t, tc.expectedMsg, msg)
		})
	}
}

func TestRabbitHandler_PollTransactionViaRabbit(t *testing.T) {
	testRabbit, err := NewTestRabbitHandler()
	require.NoError(t, err)

	defer func() {
		// close the channel and terminate the container
		testRabbit.rabbit.Channel.Close()

		if err := testRabbit.container.Terminate(testRabbit.ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()

	id, err := uuid.NewRandom()
	require.NoError(t, err)

	req := services.PollingTransactionRequest{
		TransactionId: id.String(),
	}

	rsp := services.PollingTransactionResponse{
		TransactionID:      id,
		PaydTransactionRef: gofakeit.Word(),
		Action:             gofakeit.Word(),
		Amount:             100,
		PhoneNumber:        gofakeit.Phone(),
		NetworkCode:        gofakeit.Word(),
		Naration:           gofakeit.Name(),
		PaymentStatus:      true,
		PaydUsername:       gofakeit.Name(),
	}

	rspBytes, err := json.Marshal(rsp)
	require.NoError(t, err)

	tests := []struct {
		name           string
		mockDelivery   amqp.Delivery
		expectedStatus int
		expectedMsg    gin.H
		sleep          time.Duration
	}{
		{
			name: "success",
			mockDelivery: amqp.Delivery{
				ContentType: "text/plain",
				Body:        rspBytes,
			},
			expectedStatus: http.StatusOK,
			expectedMsg:    gin.H{"response": rsp},
			sleep:          100 * time.Millisecond,
		},
		{
			name: "time out",
			mockDelivery: amqp.Delivery{
				ContentType: "text/plain",
				Body:        rspBytes,
			},
			expectedStatus: http.StatusRequestTimeout,
			expectedMsg:    gin.H{"error": "timeout waiting for response"},
			sleep:          6 * time.Second,
		},
	}

	for index, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			statusCodeChan := make(chan int, 1)
			defer close(statusCodeChan)

			msgChan := make(chan gin.H, 1)
			defer close(msgChan)

			go func(statusCodeChan chan int, msgChan chan gin.H) {
				statusCode, msg := testRabbit.rabbit.PollTransactionViaRabbit(req)
				statusCodeChan <- statusCode
				msgChan <- msg
			}(statusCodeChan, msgChan)

			// sleep so that the goroutine can send the message
			time.Sleep(tc.sleep)

			var i int
			for correlationID, responseChannel := range testRabbit.rabbit.RspMap.data {
				if i == index {
					log.Printf("CorrelationID for test case '%s': %s", tc.name, correlationID)
					tc.mockDelivery.CorrelationId = correlationID
					responseChannel <- tc.mockDelivery
				}

				i++
			}

			code := <-statusCodeChan
			msg := <-msgChan

			require.Equal(t, tc.expectedStatus, code)
			require.Equal(t, tc.expectedMsg, msg)
		})
	}
}