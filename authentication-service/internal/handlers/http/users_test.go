package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/EmilioCliff/payment-polling-app/authentication-service/internal/repository"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/pkg"
	"github.com/stretchr/testify/require"
)

var (
	TestTime = time.Date(2024, time.September, 18, 12, 0, 0, 0, time.UTC)
)

func TestHTTPServer_HandleHealthCheck(t *testing.T) {
	s := NewTestHTTPServer()

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/healthcheck", nil)
	require.NoError(t, err)

	s.server.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, `{"status":"healthy"}`, w.Body.String())
}

func mockCreateUser(user repository.User) (*repository.User, error) {
	err := user.Validate()
	if err != nil {
		return nil, pkg.Errorf(pkg.INVALID_ERROR, pkg.ErrorMessage(err))
	}

	if user.FullName == "FAIL_TEST" {
		return nil, errors.New("lets fail")
	}

	return &repository.User{
		ID: 32,
		FullName: user.FullName,
		Email: user.Email,
		Password: user.Password,
		PaydUsername: user.PaydUsername,
		PaydUsernameKey: user.PaydUsernameKey,
		PaydPasswordKey: user.PaydPasswordKey,
		PaydAccountID: user.PaydAccountID,
		CreatedAt: TestTime,
	}, nil
}

func TestHTTPServer_HandleRegisterUser(t *testing.T) {
	s := NewTestHTTPServer()

	s.UserRepository.CreateUserFunc = mockCreateUser

	req := registerUserRequest {
		FullName: "Jane",
		Email: "jane@gmail.com",
		Password: "password",
		PaydUsername: "username",
		PaydUsernameKey: "user_key",
		PaydPasswordKey: "pass_key",
		PaydAccountID: "account_id",
	}

	tests := []struct{
		name string
		payload registerUserRequest
		statusCode int
		response registerUserResponse
	}{
		{
			name: "Create User Successful",
			payload: req,
			statusCode: http.StatusOK,
			response: registerUserResponse{
				FullName: "Jane",
				Email: "jane@gmail.com",
				CreatedAt: TestTime,
			},
		},
		{
			name: "Bad Request",
			payload: registerUserRequest{},
			statusCode: http.StatusBadRequest,
			response: registerUserResponse{},
		},
		{
			name: "Internal Server Error",
			payload: registerUserRequest{
				FullName: "FAIL_TEST",
				Email: "jane@gmail.com",
				Password: "password",
				PaydUsername: "username",
				PaydUsernameKey: "user_key",
				PaydPasswordKey: "pass_key",
				PaydAccountID: "account_id",
			},
			statusCode: http.StatusInternalServerError,
			response: registerUserResponse{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			b, err := json.Marshal(tc.payload)
			require.NoError(t, err)
			req, err := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(b))
			require.NoError(t, err)

			s.server.router.ServeHTTP(w, req)
			require.Equal(t, tc.statusCode, w.Code)
			if tc.statusCode == http.StatusOK {
				expextedRsp, err := json.Marshal(tc.response)
				require.NoError(t, err)
				require.Equal(t, string(expextedRsp), w.Body.String())
			} else {
				require.Equal(t, tc.response, registerUserResponse{})
			}
		})
	}
}

func mockGetUserFunc(email string) (*repository.User, error) {
	if email == "success" {
		hashPassword, _ := pkg.GenerateHashPassword("password", 10)
		return &repository.User{
			ID: 32,
			FullName: "Jane",
			Email: "jane@gmail.com",
			Password: hashPassword,
			PaydUsername: "username",
			PaydUsernameKey: "user_key",
			PaydPasswordKey: "pass_key",
			PaydAccountID: "account_id",
			CreatedAt: time.Now(),
		}, nil
	} else if email == "unauthorized" {
		return &repository.User{
			ID: 32,
			FullName: "Jane",
			Email: "jane@gmail.com",
			Password: "UNAUTHORIZED",
			PaydUsername: "username",
			PaydUsernameKey: "user_key",
			PaydPasswordKey: "pass_key",
			PaydAccountID: "account_id",
			CreatedAt: time.Now(),
		}, nil
	}
	return nil, errors.New("user not found")
}

func TestHTTPServer_HandleLoginUser(t *testing.T) {
	s := NewTestHTTPServer()

	s.UserRepository.GetUserFunc = mockGetUserFunc

	tests := []struct{
		name string
		payload LoginUserRequest
		statusCode int
	}{
		{
			name: "Successful Log in",
			payload: LoginUserRequest{
				Email: "success",
				Password: "password",
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Unauthorized_User",
			payload: LoginUserRequest{
				Email: "unauthorized",
				Password: "password",
			},
			statusCode: http.StatusUnauthorized,
		},
		{
			name: "no user",
			payload: LoginUserRequest{
				Email: "no_user",
				Password: "password",
			},
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "bad_request",
			payload: LoginUserRequest{},
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			b, err := json.Marshal(tc.payload)
			require.NoError(t, err)
			req, err := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(b))
			require.NoError(t, err)

			s.server.router.ServeHTTP(w, req)
			require.Equal(t, tc.statusCode, w.Code)
		})
	}
}

func TestHTTPServer_ConvertPkgError(t *testing.T) {
	s := NewTestHTTPServer()

	tests := []struct{
		name string
		err string
		want int
	}{
		{
			name: "already_exists",
			err: pkg.ALREADY_EXISTS_ERROR,
			want: http.StatusConflict,
		},
		{
			name: "internal_error",
			err: pkg.INTERNAL_ERROR,
			want: http.StatusInternalServerError,
		},
		{
			name: "invalid_error",
			err: pkg.INVALID_ERROR,
			want: http.StatusBadRequest,
		},
		{
			name: "not_found",
			err: pkg.NOT_FOUND_ERROR,
			want: http.StatusNotFound,
		},
		{
			name: "not_implemented",
			err: pkg.NOT_IMPLEMENTED_ERROR,
			want: http.StatusNotImplemented,
		},
		{
			name: "authentication_error",
			err: pkg.AUTHENTICATION_ERROR,
			want: http.StatusUnauthorized,
		},
		{
			name: "default",
			err: "system_error",
			want: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := s.server.convertPkgError(tc.err)
			if got != tc.want {
				t.Errorf("convertPkgError() = %v, want %v", got, tc.want)
			}
		})
	}
}