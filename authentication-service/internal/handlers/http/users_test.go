package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/EmilioCliff/payment-polling-app/authentication-service/internal/postgres/generated"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/internal/repository"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/pkg"
	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
)

var (
	TestTime          = time.Date(2024, time.September, 18, 12, 0, 0, 0, time.UTC)
	authorizedEmail   = "login_success"
	unauthorizedEmail = "login_fail"
	notFoundEmail     = "no_user"
	defaultPassword   = "password"
)

func randomReqUser() registerUserRequest {
	return registerUserRequest{
		FullName:        gofakeit.Name(),
		Email:           gofakeit.Email(),
		Password:        gofakeit.Password(true, true, true, true, true, 7),
		PaydUsername:    gofakeit.Username(),
		PaydAccountID:   gofakeit.UUID(),
		PaydUsernameKey: gofakeit.UUID(),
		PaydPasswordKey: gofakeit.UUID(),
	}
}

func repositoryUserToGenerated(user repository.User) generated.User {
	return generated.User{
		ID:              user.ID,
		FullName:        user.FullName,
		Email:           user.Email,
		Password:        user.Password,
		PaydUsername:    user.PaydUsername,
		PaydAccountID:   user.PaydAccountID,
		PaydUsernameKey: user.PaydUsernameKey,
		PaydPasswordKey: user.PaydPasswordKey,
		CreatedAt:       user.CreatedAt,
	}
}

func TestHTTPServer_HandleHealthCheck(t *testing.T) {
	s := NewTestHTTPServer()

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/healthcheck", nil)
	require.NoError(t, err)

	s.server.router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, `{"status":"healthy"}`, w.Body.String())
}

func mockCreateUser(user repository.User) (*repository.User, error) {
	err := user.Validate()
	if err != nil {
		return nil, pkg.Errorf(pkg.INVALID_ERROR, "%v", pkg.ErrorMessage(err))
	}

	if user.FullName == "FAIL_TEST" {
		return nil, errors.New("lets fail")
	}

	return &repository.User{
		ID:              32,
		FullName:        user.FullName,
		Email:           user.Email,
		Password:        user.Password,
		PaydUsername:    user.PaydUsername,
		PaydUsernameKey: user.PaydUsernameKey,
		PaydPasswordKey: user.PaydPasswordKey,
		PaydAccountID:   user.PaydAccountID,
		CreatedAt:       TestTime,
	}, nil
}

func TestHTTPServer_HandleRegisterUser(t *testing.T) {
	s := NewTestHTTPServer()

	s.UserRepository.CreateUserFunc = mockCreateUser

	req := randomReqUser()

	tests := []struct {
		name       string
		args       registerUserRequest
		statusCode int
		rsp        registerUserResponse
	}{
		{
			name:       "Create User Successful",
			args:       req,
			statusCode: http.StatusOK,
			rsp: registerUserResponse{
				FullName:  req.FullName,
				Email:     req.Email,
				CreatedAt: TestTime,
			},
		},
		{
			name:       "Bad Request",
			args:       registerUserRequest{},
			statusCode: http.StatusBadRequest,
			rsp:        registerUserResponse{},
		},
		{
			name: "Internal Server Error",
			args: registerUserRequest{
				FullName:        "FAIL_TEST",
				Email:           req.Email,
				Password:        req.Password,
				PaydUsername:    req.PaydUsername,
				PaydUsernameKey: req.PaydUsernameKey,
				PaydPasswordKey: req.PaydPasswordKey,
				PaydAccountID:   req.PaydAccountID,
			},
			statusCode: http.StatusInternalServerError,
			rsp:        registerUserResponse{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			b, err := json.Marshal(tc.args)
			require.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(b))
			require.NoError(t, err)

			s.server.router.ServeHTTP(w, req)
			require.Equal(t, tc.statusCode, w.Code)

			if tc.statusCode == http.StatusOK {
				expextedRsp, err := json.Marshal(tc.rsp)
				require.NoError(t, err)
				require.Equal(t, string(expextedRsp), w.Body.String())
			} else {
				require.Equal(t, tc.rsp, registerUserResponse{})
			}
		})
	}
}

func mockGetUserFunc(email string) (*repository.User, error) {
	if email == authorizedEmail {
		hashPassword, _ := pkg.GenerateHashPassword(defaultPassword, 10)

		return &repository.User{
			FullName:  "Jane",
			Email:     email,
			Password:  hashPassword,
			CreatedAt: time.Now(),
		}, nil
	} else if email == unauthorizedEmail {
		return &repository.User{
			FullName:  "Jane",
			Email:     email,
			Password:  "UNAUTHORIZED",
			CreatedAt: time.Now(),
		}, nil
	}

	return nil, errors.New("user not found")
}

func TestHTTPServer_HandleLoginUser(t *testing.T) {
	s := NewTestHTTPServer()

	s.UserRepository.GetUserFunc = mockGetUserFunc

	tests := []struct {
		name       string
		payload    LoginUserRequest
		statusCode int
	}{
		{
			name: "Successful Log in",
			payload: LoginUserRequest{
				Email:    authorizedEmail,
				Password: defaultPassword,
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Unauthorized_User",
			payload: LoginUserRequest{
				Email:    unauthorizedEmail,
				Password: defaultPassword,
			},
			statusCode: http.StatusUnauthorized,
		},
		{
			name: "no user",
			payload: LoginUserRequest{
				Email:    notFoundEmail,
				Password: defaultPassword,
			},
			statusCode: http.StatusInternalServerError,
		},
		{
			name:       "bad_request",
			payload:    LoginUserRequest{},
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			b, err := json.Marshal(tc.payload)
			require.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(b))
			require.NoError(t, err)

			s.server.router.ServeHTTP(w, req)
			require.Equal(t, tc.statusCode, w.Code)
		})
	}
}

func TestHTTPServer_ConvertPkgError(t *testing.T) {
	tests := []struct {
		name string
		err  string
		want int
	}{
		{
			name: "already_exists",
			err:  pkg.ALREADY_EXISTS_ERROR,
			want: http.StatusConflict,
		},
		{
			name: "internal_error",
			err:  pkg.INTERNAL_ERROR,
			want: http.StatusInternalServerError,
		},
		{
			name: "invalid_error",
			err:  pkg.INVALID_ERROR,
			want: http.StatusBadRequest,
		},
		{
			name: "not_found",
			err:  pkg.NOT_FOUND_ERROR,
			want: http.StatusNotFound,
		},
		{
			name: "not_implemented",
			err:  pkg.NOT_IMPLEMENTED_ERROR,
			want: http.StatusNotImplemented,
		},
		{
			name: "authentication_error",
			err:  pkg.AUTHENTICATION_ERROR,
			want: http.StatusUnauthorized,
		},
		{
			name: "default",
			err:  "system_error",
			want: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := convertPkgError(tc.err)
			if got != tc.want {
				t.Errorf("convertPkgError() = %v, want %v", got, tc.want)
			}
		})
	}
}
