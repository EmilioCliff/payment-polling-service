package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/EmilioCliff/payment-polling-app/gateway-service/internal/services"
	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
)

func mockRegisterUserViagRPC(_ services.RegisterUserRequest) (int, services.RegisterUserResponse) {
	return http.StatusOK, services.RegisterUserResponse{Message: "success"}
}

func mockRegisterUserViaHttp(_ services.RegisterUserRequest) (int, services.RegisterUserResponse) {
	return http.StatusOK, services.RegisterUserResponse{Message: "success"}
}

func mockRegisterUserViaRabbit(_ services.RegisterUserRequest) (int, services.RegisterUserResponse) {
	return http.StatusOK, services.RegisterUserResponse{Message: "success"}
}

func TestHttpServer_handleRegisterUser(t *testing.T) {
	s := NewTestHttpServer()

	// change here to the communication channel you are using
	s.GrpcService.RegisterUserViagRPCFunc = mockRegisterUserViagRPC

	req := randomReq()

	tests := []struct {
		name string
		req  any
		want int
	}{
		{
			name: "success",
			req:  req,
			want: http.StatusOK,
		},
		{
			name: "missing arg",
			req:  services.LoginUserRequest{},
			want: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			b, err := json.Marshal(tc.req)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(b))
			require.NoError(t, err)

			s.server.router.ServeHTTP(w, req)
			require.Equal(t, tc.want, w.Code)
		})
	}
}

func mockLoginUserViagRPC(_ services.LoginUserRequest) (int, services.LoginUserResponse) {
	return http.StatusOK, services.LoginUserResponse{Message: "success"}
}

func mockLoginUserViaHttp(_ services.LoginUserRequest) (int, services.LoginUserResponse) {
	return http.StatusOK, services.LoginUserResponse{Message: "success"}
}

func mockLoginUserViaRabbit(_ services.LoginUserRequest) (int, services.LoginUserResponse) {
	return http.StatusOK, services.LoginUserResponse{Message: "success"}
}

func TestHttpServer_handleLoginUser(t *testing.T) {
	s := NewTestHttpServer()

	// change here to the communication channel you are using
	s.HTTPService.LoginUserViaHttpFunc = mockLoginUserViaHttp

	req := services.LoginUserRequest{
		Email:    gofakeit.Email(),
		Password: gofakeit.Password(true, true, true, true, true, 7),
	}

	tests := []struct {
		name string
		req  any
		want int
	}{
		{
			name: "success",
			req:  req,
			want: http.StatusOK,
		},
		{
			name: "missing arg",
			req:  services.RegisterUserRequest{},
			want: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			b, err := json.Marshal(tc.req)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(b))
			require.NoError(t, err)

			s.server.router.ServeHTTP(w, req)
			require.Equal(t, tc.want, w.Code)
		})
	}
}

func mockInitiatePaymentViaRabbit(_ services.InitiatePaymentRequest) (int, services.InitiatePaymentResponse) {
	return http.StatusOK, services.InitiatePaymentResponse{Message: "success"}
}

func TestHttpServer_handleInitiatePayment(t *testing.T) {
	s := NewTestHttpServer()

	accessToken, err := s.server.maker.CreateToken("user", 1, time.Minute)
	require.NoError(t, err)

	// change here to the communication channel you are using
	s.RabbitService.InitiatePaymentViaRabbitFunc = mockInitiatePaymentViaRabbit

	tests := []struct {
		name string
		req  any
		want int
	}{
		{
			name: "success",
			req: services.InitiatePaymentRequest{
				Email:       gofakeit.Email(),
				Action:      "withdrawal",
				Amount:      200,
				PhoneNumber: "phone_number",
				NetworkCode: "63902",
				Naration:    "narration",
			},
			want: http.StatusOK,
		},
		{
			name: "missing arg",
			req:  services.RegisterUserRequest{},
			want: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			b, err := json.Marshal(tc.req)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/payments/initiate", bytes.NewBuffer(b))
			require.NoError(t, err)

			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

			s.server.router.ServeHTTP(w, req)
			require.Equal(t, tc.want, w.Code)
		})
	}
}

func mockPollTransactionViaRabbit(req services.PollingTransactionRequest, userID int64) (int, services.PollingTransactionResponse) {
	if req.TransactionId == "bad_gateway" {
		return http.StatusBadGateway, services.PollingTransactionResponse{Message: "Bad gateway"}
	}

	return http.StatusOK, services.PollingTransactionResponse{Message: "success"}
}

func TestHttpServer_handlePaymentPolling(t *testing.T) {
	s := NewTestHttpServer()

	accessToken, err := s.server.maker.CreateToken("user", 1, time.Minute)
	require.NoError(t, err)

	// change here to the communication channel you are using
	s.RabbitService.PollTransactionViaRabbitFunc = mockPollTransactionViaRabbit

	tests := []struct {
		name string
		path string
		want int
	}{
		{
			name: "success",
			path: "/payments/status/123",
			want: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodGet, tc.path, nil)
			require.NoError(t, err)

			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

			s.server.router.ServeHTTP(w, req)
			require.Equal(t, tc.want, w.Code)
		})
	}
}
