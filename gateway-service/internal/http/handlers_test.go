package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EmilioCliff/payment-polling-app/gateway-service/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func mockRegisterUserViagRPC(_ services.RegisterUserRequest) (int, gin.H) {
	return http.StatusOK, gin.H{"message": "just for coverage"}
}

func mockRegisterUserViaHttp(_ services.RegisterUserRequest) (int, gin.H) {
	return http.StatusOK, gin.H{"message": "just for coverage"}
}

func mockRegisterUserViaRabbit(_ services.RegisterUserRequest) (int, gin.H) {
	return http.StatusOK, gin.H{"message": "just for coverage"}
}

func TestHttpServer_handleRegisterUser(t *testing.T) {
	s := NewTestHttpServer()

	// change here to the communication channel you are using
	s.GrpcService.RegisterUserViagRPCFunc = mockRegisterUserViagRPC

	tests := []struct{
		name string
		req any
		want int
	}{
		{
			name: "success",
			req: services.RegisterUserRequest{
				FullName: "Jane",
				Email: "jane@gmail.com",
				Password: "password",
				PaydUsername: "username",
				UsernameApiKey: "user_key",
				PasswordApiKey: "pass_key",
				PaydAccountID: "account_id",
			},
			want: http.StatusOK,
		},
		{
			name: "missing arg",
			req: services.LoginUserRequest{},
			want: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			b, err := json.Marshal(tc.req)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(b))
			require.NoError(t, err)
		
			s.server.router.ServeHTTP(w, req)
			require.Equal(t, tc.want, w.Code)
		})
	}
}


func mockLoginUserViagRPC(_ services.LoginUserRequest) (int, gin.H) {
	return http.StatusOK, gin.H{"message": "just for coverage"}
}

func mockLoginUserViaHttp(_ services.LoginUserRequest) (int, gin.H) {
	return http.StatusOK, gin.H{"message": "just for coverage"}
}

func mockLoginUserViaRabbit(_ services.LoginUserRequest) (int, gin.H) {
	return http.StatusOK, gin.H{"message": "just for coverage"}
}

func TestHttpServer_handleLoginUser(t *testing.T) {
	s := NewTestHttpServer()

	// change here to the communication channel you are using
	s.HTTPService.LoginUserViaHttpFunc = mockLoginUserViaHttp

	tests := []struct{
		name string
		req any
		want int
	}{
		{
			name: "success",
			req: services.LoginUserRequest{
				Email: "jane@gmail.com",
				Password: "password",
			},
			want: http.StatusOK,
		},
		{
			name: "missing arg",
			req: services.RegisterUserRequest{},
			want: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			b, err := json.Marshal(tc.req)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(b))
			require.NoError(t, err)
		
			s.server.router.ServeHTTP(w, req)
			require.Equal(t, tc.want, w.Code)
		})
	}
}


func mockInitiatePaymentViaRabbit(_ services.InitiatePaymentRequest) (int, gin.H) {
	return http.StatusOK, gin.H{"message": "just for coverage"}
}

func TestHttpServer_handleInitiatePayment(t *testing.T) {
	s := NewTestHttpServer()

	// change here to the communication channel you are using
	s.RabbitService.InitiatePaymentViaRabbitFunc = mockInitiatePaymentViaRabbit

	tests := []struct{
		name string
		req any
		want int
	}{
		{
			name: "success",
			req: services.InitiatePaymentRequest{
				UserID     : 1,
				Action     : "withdrawal",
				Amount     : 200,
				PhoneNumber: "phone_number",
				NetworkCode: "network",
				Naration: "narration",
			},
			want: http.StatusOK,
		},
		{
			name: "missing arg",
			req: services.RegisterUserRequest{},
			want: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			b, err := json.Marshal(tc.req)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", "/payments/initiate", bytes.NewBuffer(b))
			require.NoError(t, err)
		
			s.server.router.ServeHTTP(w, req)
			require.Equal(t, tc.want, w.Code)
		})
	}
}

func mockPollTransactionViaRabbit(req services.PollingTransactionRequest) (int, gin.H) {
	if req.TransactionId == "bad_gateway" {
		return http.StatusBadGateway, gin.H{
			"error": "Bad Gateway",
		}
	}
	return http.StatusOK, gin.H{"message": "just for coverage"}
}

func TestHttpServer_handlePaymentPolling(t *testing.T) {
	s := NewTestHttpServer()

	// change here to the communication channel you are using
	s.RabbitService.PollTransactionViaRabbitFunc = mockPollTransactionViaRabbit

	tests := []struct{
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

			req, err := http.NewRequest("GET", tc.path, nil)
			require.NoError(t, err)
		
			s.server.router.ServeHTTP(w, req)
			require.Equal(t, tc.want, w.Code)
		})
	}
}