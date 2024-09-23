package Grpc

import (
	"context"
	"database/sql"
	"errors"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/EmilioCliff/payment-polling-app/authentication-service/internal/repository"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/pkg"
	"github.com/EmilioCliff/payment-polling-service/shared-grpc/pb"
	"github.com/brianvoe/gofakeit"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	TestTime                = time.Date(2024, time.September, 18, 12, 0, 0, 0, time.UTC)
	authorizedEmail         = "login_success"
	unauthorizedEmail       = "login_fail"
	notFoundEmail           = "no_user"
	foundID           int64 = 32
	notFoundID        int64 = 33
	errorID           int64 = 34
)

func mockCreateUser(user repository.User) (*repository.User, error) {
	err := user.Validate()
	if err != nil {
		return nil, pkg.Errorf(pkg.INVALID_ERROR, "%v", pkg.ErrorMessage(err))
	}

	if user.FullName == "FAIL_TEST" {
		return nil, errors.New("lets fail")
	}

	user.CreatedAt = TestTime

	return &user, nil
}

func TestGRPCServer_RegisterUser(t *testing.T) {
	s := NewTestGRPCServer()

	s.UserRepository.CreateUserFunc = mockCreateUser

	user := randomPbUser()

	tests := []struct {
		name    string
		args    *pb.RegisterUserRequest
		want    *pb.RegisterUserResponse
		wantErr bool
	}{
		{
			name: "Create User Successful",
			args: user,
			want: &pb.RegisterUserResponse{
				Fullname:  user.Fullname,
				Email:     user.Email,
				CreatedAt: timestamppb.New(TestTime),
			},
			wantErr: false,
		},
		{
			name:    "Missing values",
			args:    &pb.RegisterUserRequest{},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.server.RegisterUser(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("register user error = %v, wantErr %v", err, tt.wantErr)
			}

			if !proto.Equal(got, tt.want) {
				t.Errorf("register user got = %v, want %v", got, tt.want)
			}
		})
	}
}

func mockGetUserByID(id int64) (*repository.User, error) {
	if id == foundID {
		return &repository.User{
			ID:              foundID,
			FullName:        gofakeit.Name(),
			Email:           gofakeit.Email(),
			Password:        gofakeit.Password(true, true, true, true, true, 7),
			PaydUsername:    "username",
			PaydUsernameKey: "user_key",
			PaydPasswordKey: "pass_key",
			PaydAccountID:   "account_id",
			CreatedAt:       TestTime,
		}, nil
	} else if id == notFoundID {
		return nil, sql.ErrNoRows
	}

	return nil, errors.New("internal error")
}

func TestGRPCServer_GetUser(t *testing.T) {
	s := NewTestGRPCServer()

	s.UserRepository.GetUserByIDFunc = mockGetUserByID

	tests := []struct {
		name    string
		args    *pb.GetUserRequest
		want    *pb.GetUserResponse
		wantErr bool
	}{
		{
			name: "User found",
			args: &pb.GetUserRequest{UserId: foundID},
			want: &pb.GetUserResponse{
				PaydUsername:    "username",
				PaydUsernameKey: "user_key",
				PaydPasswordKey: "pass_key",
				PaydAccountId:   "account_id",
			},
			wantErr: false,
		},
		{
			name:    "User not found",
			args:    &pb.GetUserRequest{UserId: notFoundID},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "internal error",
			args:    &pb.GetUserRequest{UserId: errorID},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := s.server.GetUser(context.Background(), tc.args)
			if (err != nil) != tc.wantErr {
				t.Errorf("get user error = %v, wantErr %v", err, tc.wantErr)
			}

			if !proto.Equal(got, tc.want) {
				t.Errorf("login user got = %v, want %v", got, tc.want)
			}
		})
	}
}

func mockGetUserFunc(email string) (*repository.User, error) {
	if email == authorizedEmail {
		hashPassword, _ := pkg.GenerateHashPassword("password", 10)

		rsp := randomRepoUser()
		rsp.Password = hashPassword

		return rsp, nil
	} else if email == unauthorizedEmail {
		rsp := randomRepoUser()
		rsp.Password = "UNAUTHORIZED"

		return rsp, nil
	}

	return nil, errors.New("user not found")
}

func TestGRPCServer_LoginUser(t *testing.T) {
	s := NewTestGRPCServer()

	s.UserRepository.GetUserFunc = mockGetUserFunc

	tests := []struct {
		name    string
		args    *pb.LoginUserRequest
		want    *pb.LoginUserResponse
		wantErr bool
	}{
		{
			name:    "Successful Log in",
			args:    &pb.LoginUserRequest{Email: authorizedEmail, Password: "password"},
			want:    &pb.LoginUserResponse{AccessToken: "token"},
			wantErr: false,
		},
		{
			name:    "Wrong password",
			args:    &pb.LoginUserRequest{Email: unauthorizedEmail, Password: "password"},
			want:    &pb.LoginUserResponse{},
			wantErr: true,
		},
		{
			name:    "No user found",
			args:    &pb.LoginUserRequest{Email: notFoundEmail, Password: "password"},
			want:    &pb.LoginUserResponse{},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := s.server.LoginUser(context.Background(), tc.args)
			if (err != nil) != tc.wantErr {
				t.Errorf("login user error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestGRPCServer_convertPkgError(t *testing.T) {
	tests := []struct {
		name string
		err  string
		want codes.Code
	}{
		{
			name: "already_exists",
			err:  pkg.ALREADY_EXISTS_ERROR,
			want: codes.AlreadyExists,
		},
		{
			name: "internal_error",
			err:  pkg.INTERNAL_ERROR,
			want: codes.Internal,
		},
		{
			name: "invalid_error",
			err:  pkg.INVALID_ERROR,
			want: codes.InvalidArgument,
		},
		{
			name: "not_found",
			err:  pkg.NOT_FOUND_ERROR,
			want: codes.NotFound,
		},
		{
			name: "not_implemented",
			err:  pkg.NOT_IMPLEMENTED_ERROR,
			want: codes.Unimplemented,
		},
		{
			name: "authentication_error",
			err:  pkg.AUTHENTICATION_ERROR,
			want: codes.Unauthenticated,
		},
		{
			name: "default",
			err:  "system_error",
			want: codes.Internal,
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

func randomPbUser() *pb.RegisterUserRequest {
	return &pb.RegisterUserRequest{
		Fullname:           gofakeit.Name(),
		Email:              gofakeit.Email(),
		Password:           gofakeit.Password(true, true, true, true, true, 7),
		PaydUsername:       gofakeit.Username(),
		PaydAccountId:      gofakeit.UUID(),
		PaydUsernameApiKey: gofakeit.UUID(),
		PaydPasswordApiKey: gofakeit.UUID(),
	}
}

func randomRepoUser() *repository.User {
	return &repository.User{
		ID:              int64(rand.IntN(100)),
		FullName:        gofakeit.Name(),
		Email:           gofakeit.Email(),
		Password:        gofakeit.Password(true, true, true, true, true, 7),
		PaydUsername:    gofakeit.Username(),
		PaydAccountID:   gofakeit.UUID(),
		PaydUsernameKey: gofakeit.UUID(),
		PaydPasswordKey: gofakeit.UUID(),
		CreatedAt:       TestTime,
	}
}
