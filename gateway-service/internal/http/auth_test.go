package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/EmilioCliff/payment-polling-app/gateway-service/internal/services"
	"github.com/EmilioCliff/payment-polling-app/gateway-service/pkg"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

var (
	TestTime = time.Date(2024, time.September, 18, 12, 0, 0, 0, time.UTC)
)

func TestHTTPService_RegisterUserViaHttp(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		if r.URL.Path == "/auth/register" && r.Method == http.MethodPost {
			var req services.RegisterUserRequest

			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(services.RegisterUserResponse{
					Message: "fail",
					StatusCode: http.StatusInternalServerError,
				})

				return
			}

			if req.Email == "" {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(services.RegisterUserResponse{
					Message: "missing values",
					StatusCode: http.StatusBadRequest,
				})

				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(services.RegisterUserResponse{
				FullName: req.FullName,
				Email: req.Email,
				CreatedAt: TestTime,
			})

		} else {
			http.NotFound(w, r)
		}
	}))
	defer testServer.Close()

	s := NewHTTPService(pkg.Config{
		AUTH_HTTP_PORT: strings.TrimPrefix(testServer.URL, "http://"),
	})

	tests := []struct{
		name string
		req services.RegisterUserRequest
		wantStatusCode int
		wantRsp gin.H
	}{
		{
			name: "success",
			req: services.RegisterUserRequest{
				FullName: "Jane",
				Email: "jane@gmail.com",
				Password: "secret",
				PaydUsername: "username",
				UsernameApiKey: "user_key",
				PasswordApiKey: "pass_key",
				PaydAccountID: "account_id",
			},
			wantStatusCode: http.StatusOK,
			wantRsp: gin.H{"response": services.RegisterUserResponse{
				FullName: "Jane",
				Email: "jane@gmail.com",
				CreatedAt: TestTime,
			}},
		},
		{
			name: "failed request",
			req: services.RegisterUserRequest{},
			wantStatusCode: http.StatusBadRequest,
			wantRsp: gin.H{"message": "Failed to register new user", "error": "missing values"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			statusCode, msg := s.RegisterUserViaHttp(tc.req)
			require.Equal(t, statusCode, tc.wantStatusCode)
			require.Equal(t, msg, tc.wantRsp)
		})
	}
}

func TestHTTPServer_LoginUserViaHttp(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.URL.Path == "/auth/login" && r.Method == http.MethodPost {
			var req services.LoginUserRequest

			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(services.LoginUserResponse{
					Message: "fail",
					StatusCode: http.StatusInternalServerError,
				})

				return
			}

			if req.Password == "test_fail" {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(services.LoginUserResponse{
					Message: "fail",
					StatusCode: http.StatusUnauthorized,
				})
			} else {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(services.LoginUserResponse{
					AccessToken: "token",
					FullName: "jane doe",
					Email:  req.Email,
					CreatedAt: TestTime,
				})
			}
		} else {
			http.NotFound(w, r)
		}
	}))
	defer testServer.Close()

	s := NewHTTPService(pkg.Config{
		AUTH_HTTP_PORT: strings.TrimPrefix(testServer.URL, "http://"),
	})

	tests := []struct{
		name string
		req services.LoginUserRequest
		wantStatusCode int
		wantRsp gin.H
	}{
		{
			name: "success",
			req: services.LoginUserRequest{
				Email: "jane@gmail.com",
				Password: "secret",
			},
			wantStatusCode: http.StatusOK,
			wantRsp: gin.H{"response": services.LoginUserResponse{
				AccessToken: "token",
				FullName: "jane doe",
				Email: "jane@gmail.com",
				CreatedAt: TestTime,
			}},
		},
		{
			name: "failed request",
			req: services.LoginUserRequest{
				Email: "jane@gmail.com",
				Password: "test_fail",
			},
			wantStatusCode: http.StatusUnauthorized,
			wantRsp: gin.H{"message": "Failed to login user", "error": "fail"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			statusCode, msg := s.LoginUserViaHttp(tc.req)
			require.Equal(t, statusCode, tc.wantStatusCode)
			require.Equal(t, msg, tc.wantRsp)
		})
	}
}