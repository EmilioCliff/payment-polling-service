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
	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
)

var TestTime = time.Date(2024, time.September, 18, 12, 0, 0, 0, time.UTC)

func TestHTTPService_RegisterUserViaHttp(t *testing.T) {
	// creates test server to handle register request.
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.URL.Path == "/auth/register" && r.Method == http.MethodPost {
			var req services.RegisterUserRequest

			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(services.RegisterUserResponse{
					Message:    "fail",
					StatusCode: http.StatusInternalServerError,
				})

				return
			}

			if req.Email == "" {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(services.RegisterUserResponse{
					Message:    "missing values",
					StatusCode: http.StatusBadRequest,
				})

				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(services.RegisterUserResponse{
				FullName:  req.FullName,
				Email:     req.Email,
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

	req := randomReq()

	tests := []struct {
		name           string
		req            services.RegisterUserRequest
		wantStatusCode int
		wantRsp        services.RegisterUserResponse
	}{
		{
			name:           "success",
			req:            req,
			wantStatusCode: http.StatusOK,
			wantRsp: services.RegisterUserResponse{
				FullName:  req.FullName,
				Email:     req.Email,
				CreatedAt: TestTime,
			},
		},
		{
			name:           "failed request",
			req:            services.RegisterUserRequest{},
			wantStatusCode: http.StatusBadRequest,
			wantRsp:        services.RegisterUserResponse{Message: "missing values", StatusCode: http.StatusBadRequest},
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
	// creates test server to handle login request
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.URL.Path == "/auth/login" && r.Method == http.MethodPost {
			var req services.LoginUserRequest

			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(services.LoginUserResponse{
					Message:    "fail",
					StatusCode: http.StatusInternalServerError,
				})

				return
			}

			if req.Password == "test_fail" {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(services.LoginUserResponse{
					Message:    "fail",
					StatusCode: http.StatusUnauthorized,
				})
			} else {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(services.LoginUserResponse{
					AccessToken: "token",
					FullName:    "jane doe",
					Email:       req.Email,
					CreatedAt:   TestTime,
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

	req := services.LoginUserRequest{
		Email:    gofakeit.Email(),
		Password: gofakeit.Password(true, true, true, true, true, 7),
	}

	tests := []struct {
		name           string
		req            services.LoginUserRequest
		wantStatusCode int
		wantRsp        services.LoginUserResponse
	}{
		{
			name:           "success",
			req:            req,
			wantStatusCode: http.StatusOK,
			wantRsp: services.LoginUserResponse{
				AccessToken: "token",
				FullName:    "jane doe",
				Email:       req.Email,
				CreatedAt:   TestTime},
		},
		{
			name: "failed request",
			req: services.LoginUserRequest{
				Email:    req.Email,
				Password: "test_fail",
			},
			wantStatusCode: http.StatusUnauthorized,
			wantRsp:        services.LoginUserResponse{Message: "fail", StatusCode: http.StatusUnauthorized},
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
