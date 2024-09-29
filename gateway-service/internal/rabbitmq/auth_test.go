package rabbitmq

import (
	"encoding/json"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/EmilioCliff/payment-polling-app/gateway-service/internal/services"
	"github.com/EmilioCliff/payment-polling-app/gateway-service/pkg"
	"github.com/brianvoe/gofakeit"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/require"
)

var TestTime = time.Date(2024, time.September, 18, 12, 0, 0, 0, time.UTC)

func TestRabbitHandler_RegisterUserViaRabbit(t *testing.T) {
	pkg.SkipCI(t)

	testRabbit, err := NewTestRabbitHandler()
	require.NoError(t, err)

	defer func() {
		// close the channel and terminate the container
		testRabbit.rabbit.Channel.Close()

		if err := testRabbit.container.Terminate(testRabbit.ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()

	req := randomReq()

	rsp := services.RegisterUserResponse{
		FullName:  req.FullName,
		Email:     req.Email,
		CreatedAt: TestTime,
	}

	rspBytes, err := json.Marshal(rsp)
	require.NoError(t, err)

	tests := []struct {
		name           string
		mockDelivery   amqp.Delivery
		expectedStatus int
		expectedMsg    services.RegisterUserResponse
		sleep          time.Duration
	}{
		{
			name: "success",
			mockDelivery: amqp.Delivery{
				ContentType: "text/plain",
				Body:        rspBytes,
			},
			expectedStatus: http.StatusOK,
			expectedMsg:    rsp,
			sleep:          100 * time.Millisecond,
		},
		{
			name: "time out",
			mockDelivery: amqp.Delivery{
				ContentType: "text/plain",
				Body:        rspBytes,
			},
			expectedStatus: http.StatusRequestTimeout,
			expectedMsg: services.RegisterUserResponse{
				Message:    "timeout waiting for response. Try again",
				StatusCode: http.StatusInternalServerError,
			},
			sleep: 6 * time.Second,
		},
	}

	for index, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			statusCodeChan := make(chan int, 1)
			defer close(statusCodeChan)

			msgChan := make(chan services.RegisterUserResponse, 1)
			defer close(msgChan)

			go func(statusCodeChan chan int, msgChan chan services.RegisterUserResponse) {
				statusCode, msg := testRabbit.rabbit.RegisterUserViaRabbit(req)
				statusCodeChan <- statusCode
				msgChan <- msg
			}(statusCodeChan, msgChan)

			// sleep so that the goroutine can send the message
			time.Sleep(tc.sleep)

			var i int

			testRabbit.rabbit.RspMap.mu.RLock()
			for correlationID, responseChannel := range testRabbit.rabbit.RspMap.data {
				if i == index {
					log.Printf("CorrelationID for test case '%s': %s", tc.name, correlationID)
					tc.mockDelivery.CorrelationId = correlationID
					responseChannel <- tc.mockDelivery
				}

				i++
			}
			testRabbit.rabbit.RspMap.mu.RUnlock()

			code := <-statusCodeChan
			msg := <-msgChan

			require.Equal(t, tc.expectedStatus, code)
			require.Equal(t, tc.expectedMsg, msg)
		})
	}
}

func TestRabbitHandler_LoginUserViaRabbit(t *testing.T) {
	testRabbit, err := NewTestRabbitHandler()
	require.NoError(t, err)

	defer func() {
		// close the channel and terminate the container
		testRabbit.rabbit.Channel.Close()

		if err := testRabbit.container.Terminate(testRabbit.ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()

	req := services.LoginUserRequest{
		Email:    gofakeit.Email(),
		Password: gofakeit.Password(true, true, true, true, true, 7),
	}

	rsp := services.LoginUserResponse{
		Email:       req.Email,
		FullName:    gofakeit.Name(),
		AccessToken: gofakeit.UUID(),
		CreatedAt:   TestTime,
	}

	rspBytes, err := json.Marshal(rsp)
	require.NoError(t, err)

	tests := []struct {
		name           string
		mockDelivery   amqp.Delivery
		expectedStatus int
		expectedMsg    services.LoginUserResponse
		sleep          time.Duration
	}{
		{
			name: "success",
			mockDelivery: amqp.Delivery{
				ContentType: "text/plain",
				Body:        rspBytes,
			},
			expectedStatus: http.StatusOK,
			expectedMsg:    rsp,
			sleep:          100 * time.Millisecond,
		},
		{
			name: "time out",
			mockDelivery: amqp.Delivery{
				ContentType: "text/plain",
				Body:        rspBytes,
			},
			expectedStatus: http.StatusRequestTimeout,
			expectedMsg: services.LoginUserResponse{
				Message:    "timeout waiting for response. Try again",
				StatusCode: http.StatusInternalServerError,
			},
			sleep: 6 * time.Second,
		},
	}

	for index, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			statusCodeChan := make(chan int, 1)
			defer close(statusCodeChan)

			msgChan := make(chan services.LoginUserResponse, 1)
			defer close(msgChan)

			go func(statusCodeChan chan int, msgChan chan services.LoginUserResponse) {
				statusCode, msg := testRabbit.rabbit.LoginUserViaRabbit(req)
				statusCodeChan <- statusCode
				msgChan <- msg
			}(statusCodeChan, msgChan)

			// sleep so that the goroutine can send the message
			time.Sleep(tc.sleep)

			var i int

			testRabbit.rabbit.RspMap.mu.RLock()
			for correlationID, responseChannel := range testRabbit.rabbit.RspMap.data {
				if i == index {
					log.Printf("CorrelationID for test case '%s': %s", tc.name, correlationID)
					tc.mockDelivery.CorrelationId = correlationID
					responseChannel <- tc.mockDelivery
				}

				i++
			}
			testRabbit.rabbit.RspMap.mu.RUnlock()

			code := <-statusCodeChan
			msg := <-msgChan

			require.Equal(t, tc.expectedStatus, code)
			require.Equal(t, tc.expectedMsg, msg)
		})
	}
}

func randomReq() services.RegisterUserRequest {
	return services.RegisterUserRequest{
		FullName:       gofakeit.Name(),
		Email:          gofakeit.Email(),
		Password:       gofakeit.Password(true, true, true, true, true, 7),
		PaydUsername:   gofakeit.Username(),
		PaydAccountID:  gofakeit.UUID(),
		UsernameApiKey: gofakeit.UUID(),
		PasswordApiKey: gofakeit.UUID(),
	}
}
