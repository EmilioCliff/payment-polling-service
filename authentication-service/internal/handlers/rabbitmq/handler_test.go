package rabbitmq_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/EmilioCliff/payment-polling-app/authentication-service/internal/handlers/rabbitmq"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/internal/repository"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/pkg"
	"golang.org/x/crypto/bcrypt"
)

var TestTime = time.Date(2024, time.September, 18, 12, 0, 0, 0, time.UTC)

func mockCreateUser(user repository.User) (*repository.User, error) {
	err := user.Validate()
	if err != nil {
		return nil, err
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

func TestRabbitConn_HandleRegisterUser(t *testing.T) {
	r := NewTestRabbitConn()

	r.UserRepository.CreateUserFunc = mockCreateUser

	tests := []struct {
		name    string
		request rabbitmq.RegisterUserRequest
		want    any
		wantErr bool
	}{
		{
			name: "Create User Successful",
			request: rabbitmq.RegisterUserRequest{
				FullName:        "Jane",
				Email:           "jane@gmail.com",
				Password:        "password",
				PaydUsername:    "username",
				PaydUsernameKey: "user_key",
				PaydPasswordKey: "pass_key",
				PaydAccountID:   "account_id",
			},
			want: rabbitmq.RegisterUserResponse{
				FullName:  "Jane",
				Email:     "jane@gmail.com",
				CreatedAt: TestTime,
			},
			wantErr: false,
		},
		{
			name:    "Missing values",
			request: rabbitmq.RegisterUserRequest{},
			want: rabbitmq.ErrorResponse{
				StatusCode: 400,
				Message:    "failed to create user: full_name is required",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBytes := r.rabbitConn.HandleRegisterUser(tt.request)

			if tt.wantErr {
				var gotErrResp rabbitmq.ErrorResponse
				if err := json.Unmarshal(gotBytes, &gotErrResp); err != nil {
					t.Fatalf("failed to unmarshal error response: %v", err)
				}

				if gotErrResp != tt.want {
					t.Errorf("got error response %v, want %v", gotErrResp, tt.want)
				}
			} else {
				var gotSuccessResp rabbitmq.RegisterUserResponse
				if err := json.Unmarshal(gotBytes, &gotSuccessResp); err != nil {
					t.Fatalf("failed to unmarshal success response: %v", err)
				}

				if gotSuccessResp != tt.want {
					t.Errorf("got success response %v, want %v", gotSuccessResp, tt.want)
				}
			}
		})
	}
}

func mockGetUserFunc(email string) (*repository.User, error) {
	if email == "success" {
		hashPassword, _ := pkg.GenerateHashPassword("password", 10)

		return &repository.User{
			ID:              32,
			FullName:        "Jane",
			Email:           "jane@gmail.com",
			Password:        hashPassword,
			PaydUsername:    "username",
			PaydUsernameKey: "user_key",
			PaydPasswordKey: "pass_key",
			PaydAccountID:   "account_id",
			CreatedAt:       TestTime,
		}, nil
	} else if email == "unauthorized" {
		hashPassword, _ := pkg.GenerateHashPassword("differen_hash", 10)

		return &repository.User{
			ID:              32,
			FullName:        "Jane",
			Email:           "jane@gmail.com",
			Password:        hashPassword,
			PaydUsername:    "username",
			PaydUsernameKey: "user_key",
			PaydPasswordKey: "pass_key",
			PaydAccountID:   "account_id",
			CreatedAt:       TestTime,
		}, nil
	}

	return nil, pkg.Errorf(pkg.NOT_FOUND_ERROR, "error getting user: %s", email)
}

func TestRabbitConn_HandleLoginUser(t *testing.T) {
	r := NewTestRabbitConn()

	r.UserRepository.GetUserFunc = mockGetUserFunc

	accessToken, _ := r.rabbitConn.Maker.CreateToken(
		"jane@gmail.com",
		r.rabbitConn.Config.TOKEN_DURATION,
	)

	tests := []struct {
		name    string
		args    rabbitmq.LoginUserRequest
		want    any
		wantErr bool
	}{
		{
			name: "Successful Log in",
			args: rabbitmq.LoginUserRequest{Email: "success", Password: "password"},
			want: rabbitmq.LoginUserResponse{
				Email:       "jane@gmail.com",
				AccessToken: accessToken,
				FullName:    "Jane",
				CreatedAt:   TestTime,
			},
			wantErr: false,
		},
		{
			name: "Wrong password",
			args: rabbitmq.LoginUserRequest{Email: "unauthorized", Password: "password"},
			want: rabbitmq.ErrorResponse{
				StatusCode: 401,
				Message: fmt.Sprintf(
					"Error comparing passwords: %v",
					bcrypt.ErrMismatchedHashAndPassword,
				),
			},
			wantErr: true,
		},
		{
			name: "No user found",
			args: rabbitmq.LoginUserRequest{Email: "no_user", Password: "password"},
			want: rabbitmq.ErrorResponse{
				StatusCode: 404,
				Message:    "failed to login user: error getting user: no_user",
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotBytes := r.rabbitConn.HandleLoginUser(tc.args)

			if tc.wantErr {
				var gotErrResp rabbitmq.ErrorResponse
				if err := json.Unmarshal(gotBytes, &gotErrResp); err != nil {
					t.Fatalf("failed to unmarshal error response: %v", err)
				}

				if gotErrResp != tc.want {
					t.Errorf("got error response %v, want %v", gotErrResp, tc.want)
				}
			} else {
				var gotSuccessResp rabbitmq.LoginUserResponse
				if err := json.Unmarshal(gotBytes, &gotSuccessResp); err != nil {
					t.Fatalf("failed to unmarshal success response: %v", err)
				}

				if gotSuccessResp.Email != "jane@gmail.com" {
					t.Errorf("got success response %v, want %v", gotSuccessResp.AccessToken, accessToken)
				}
			}
		})
	}
}
