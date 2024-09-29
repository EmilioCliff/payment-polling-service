package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EmilioCliff/payment-polling-app/payment-service/internal/mock"
	"github.com/EmilioCliff/payment-polling-app/payment-service/internal/repository"
	"github.com/EmilioCliff/payment-polling-app/payment-service/pkg"
	"github.com/brianvoe/gofakeit"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type TestHttpServer struct {
	server *HttpServer

	TransactionRepository mock.MockTransactionRepository
}

func NewTestHttpServer() *TestHttpServer {
	gin.SetMode(gin.TestMode)

	s := &TestHttpServer{
		server: NewHttpServer(),
	}

	s.server.TransactionRepository = &s.TransactionRepository

	return s
}

func TestHTTPServer_HandleHealthCheck(t *testing.T) {
	s := NewTestHttpServer()

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/healthcheck", nil)
	require.NoError(t, err)

	s.server.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, `{"status":"healthy"}`, w.Body.String())
}

func mockUpdateTransactionFunc(ctx context.Context, id uuid.UUID, req repository.TransactionUpdate) (*repository.Transaction, error) {
	if id == uuid.Nil {
		return nil, pkg.Errorf(pkg.INVALID_ERROR, "id cannot be empty")
	}

	return nil, nil
}

func TestHttpServer_handleCallback(t *testing.T) {
	s := NewTestHttpServer()

	s.TransactionRepository.UpdateTransactionFunc = mockUpdateTransactionFunc

	tests := []struct {
		name string
		id   string
		req  map[string]interface{}
		want int
	}{
		{
			name: "success",
			id:   gofakeit.UUID(),
			req: map[string]interface{}{
				"amount":                20,
				"forward_url":           "https://603f-105-163-2-208.ngrok-free.app/callback",
				"order_id":              "",
				"phone_number":          "254718750145",
				"remarks":               "Transaction processed successfully with reference: [SIR041D2V4]",
				"result_code":           0,
				"third_party_trans_id":  "SIR041D2V4",
				"transaction_date":      "November",
				"transaction_reference": "RIB181154135",
			},
			want: http.StatusOK,
		},
		{
			name: "invalid uuid",
			id:   "invalid uuid",
			req:  map[string]interface{}{},
			want: http.StatusBadRequest,
		},
		{
			name: "nil uuid",
			id:   uuid.Nil.String(),
			req:  nil,
			want: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			b, err := json.Marshal(tc.req)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/transaction/%v", tc.id), bytes.NewReader(b))
			require.NoError(t, err)

			s.server.router.ServeHTTP(w, req)

			require.Equal(t, tc.want, w.Code)
		})
	}
}
