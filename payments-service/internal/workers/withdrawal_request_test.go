package workers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EmilioCliff/payment-polling-app/payment-service/internal/services"
	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/require"
)

func TestRedisTaskProcessor_ProcessWithdrawalRequestTask(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodPost {
			var req map[string]interface{}

			err := json.NewDecoder(r.Body).Decode(&req)
			require.NoError(t, err)

			if req["phone_number"] == "test_fail" {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"error_message": "test failing",
				})

				return
			}

			w.WriteHeader(http.StatusAccepted)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"merchantRequestID": gofakeit.UUID(),
				"message":           gofakeit.Sentence(10),
			})
		} else {
			http.NotFound(w, r)
		}
	}))
	defer testServer.Close()

	p := NewTestRedisProcessor()

	p.TransactionRepository.CreateTransactionFunc = mockCreateTransactionFunc

	p.redisProcessor.withdrawalUrl = testServer.URL

	tests := []struct {
		name    string
		req     services.SendPaymentWithdrawalRequestPayload
		wantErr bool
	}{
		{
			name: "success",
			req: services.SendPaymentWithdrawalRequestPayload{
				TransactionID:      uuid.New(),
				UserID:             1,
				Action:             "withdraw",
				Amount:             100,
				PhoneNumber:        gofakeit.Phone(),
				NetworkCode:        "63902",
				Naration:           gofakeit.Sentence(10),
				PaydUsername:       gofakeit.Name(),
				PaydAccountID:      gofakeit.UUID(),
				PaydPasswordApiKey: gofakeit.UUID(),
				PaydUsernameApiKey: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "fail request",
			req: services.SendPaymentWithdrawalRequestPayload{
				TransactionID:      uuid.New(),
				UserID:             32,
				Action:             "withdraw",
				Amount:             100,
				PhoneNumber:        "test_fail",
				NetworkCode:        "63902",
				Naration:           gofakeit.Sentence(10),
				PaydUsername:       gofakeit.Name(),
				PaydAccountID:      gofakeit.UUID(),
				PaydPasswordApiKey: gofakeit.UUID(),
				PaydUsernameApiKey: gofakeit.UUID(),
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			payloadBytes, err := json.Marshal(tc.req)
			require.NoError(t, err)

			payload := asynq.NewTask("payment_request", payloadBytes)

			err = p.redisProcessor.ProcessWithdrawalRequestTask(context.Background(), payload)
			if (err != nil) != tc.wantErr {
				t.Errorf("ProcessPaymentRequestTask() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
