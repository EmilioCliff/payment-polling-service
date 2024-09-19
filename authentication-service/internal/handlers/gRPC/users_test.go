package gRPC

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/EmilioCliff/payment-polling-app/authentication-service/internal/repository"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/pkg"
	pb "github.com/EmilioCliff/payment-polling-service/shared-grpc/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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
		CreatedAt: time.Date(2024, time.September, 18, 12, 0, 0, 0, time.UTC),
	}, nil
}

func TestGRPCServer_RegisterUser(t *testing.T) {
	s := NewTestGRPCServer()

	s.UserRepository.CreateUserFunc = mockCreateUser

	tests := []struct{
		name string
		request *pb.RegisterUserRequest
		want *pb.RegisterUserResponse
		wantErr bool
	}{
		{
			name: "Create User Successful",
			request: &pb.RegisterUserRequest{
				Fullname: "Jane",
				Email: "jane@gmail.com",
				Password: "password",
				PaydUsername: "username",
				PaydUsernameApiKey: "user_key",
				PaydPasswordApiKey: "pass_key",
				PaydAccountId: "account_id",
			},
			want: &pb.RegisterUserResponse{
				Fullname: "Jane",
				Email: "jane@gmail.com",
				CreatedAt: timestamppb.New(time.Date(2024, time.September, 18, 12, 0, 0, 0, time.UTC)),
			},
			wantErr: false,
		},
		{
			name: "Missing values",
			request: &pb.RegisterUserRequest{},
			want: nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := s.server.RegisterUser(context.Background(), tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("register user error = %v, wantErr %v", err, tt.wantErr)
			}
			if  !proto.Equal(got, tt.want) {
				t.Errorf("register user got = %v, want %v", got, tt.want)
			}
        })
    }
}

func mockGetUserByID(id int64) (*repository.User, error) {
	if id == 32 {
		return &repository.User{
			ID: 32,
			FullName: "Jane",
			Email: "jane@gmail.com",
			Password: "password",
			PaydUsername: "username",
			PaydUsernameKey: "user_key",
			PaydPasswordKey: "pass_key",
			PaydAccountID: "account_id",
			CreatedAt: time.Date(2024, time.September, 18, 12, 0, 0, 0, time.UTC),
		}, nil
	} else if id == 33 {
		return nil, sql.ErrNoRows
	}
	return nil, errors.New("user not found")
}

func TestGRPCServer_GetUser(t *testing.T) {
	s := NewTestGRPCServer()

	s.UserRepository.GetUserByIDFunc = mockGetUserByID

	tests := []struct{
		name string
		id *pb.GetUserRequest
		want *pb.GetUserResponse
		wantErr bool
	}{
		{
			name: "User found",
			id: &pb.GetUserRequest{
				UserId: 32,
			},
			want: &pb.GetUserResponse{
				PaydUsername: "username",
				PaydUsernameKey: "user_key",
				PaydPasswordKey: "pass_key",
				PaydAccountId: "account_id",
			},
			wantErr: false,
		},
		{
			name: "User not found",
			id: &pb.GetUserRequest{
				UserId: 33,
			},
			want: nil,
			wantErr: true,
		},
		{
			name: "internal error",
			id: &pb.GetUserRequest{
				UserId: 0,
			},
			want: nil,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got ,err := s.server.GetUser(context.Background(), tc.id)
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

func TestGRPCServer_LoginUser(t *testing.T) {
	s := NewTestGRPCServer()

	s.UserRepository.GetUserFunc = mockGetUserFunc

	tests := []struct{
		name string
		args *pb.LoginUserRequest
		want *pb.LoginUserResponse
		wantErr bool
	}{
		{
			name: "Successful Log in",
			args: &pb.LoginUserRequest{
				Email: "success",
				Password: "password",
			},
			want: &pb.LoginUserResponse{
				AccessToken: "token",
			},
			wantErr: false,
		},
		{
			name: "Wrong password",
			args: &pb.LoginUserRequest{
				Email: "unauthorized",
				Password: "password",
			},
			want: &pb.LoginUserResponse{},
			wantErr: true,
		},
		{
			name: "No user found",
			args: &pb.LoginUserRequest{
				Email: "no user found",
				Password: "password",
			},
			want: &pb.LoginUserResponse{},
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

func TestGRPCServer_convertPkgError (t *testing.T) {
	s := NewTestGRPCServer()

	tests := []struct{
		name string
		err string
		want codes.Code
	}{
		{
			name: "already_exists",
			err: pkg.ALREADY_EXISTS_ERROR,
			want: codes.AlreadyExists,
		},
		{
			name: "internal_error",
			err: pkg.INTERNAL_ERROR,
			want: codes.Internal,
		},
		{
			name: "invalid_error",
			err: pkg.INVALID_ERROR,
			want: codes.InvalidArgument,
		},
		{
			name: "not_found",
			err: pkg.NOT_FOUND_ERROR,
			want: codes.NotFound,
		},
		{
			name: "not_implemented",
			err: pkg.NOT_IMPLEMENTED_ERROR,
			want: codes.Unimplemented,
		},
		{
			name: "authentication_error",
			err: pkg.AUTHENTICATION_ERROR,
			want: codes.Unauthenticated,
		},
		{
			name: "default",
			err: "system_error",
			want: codes.Internal,
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