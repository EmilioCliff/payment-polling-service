package gRPC

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/EmilioCliff/payment-polling-app/gateway-service/internal/services"
	grpcMock "github.com/EmilioCliff/payment-polling-service/shared-grpc/mock"
	"github.com/EmilioCliff/payment-polling-service/shared-grpc/pb"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	TestTime = time.Date(2024, time.September, 18, 12, 0, 0, 0, time.UTC)
)

func TestGrpcClient_RegisterUserViagRPC(t *testing.T) {
	g := NewTestGrpcClient()

	ctrl := gomock.NewController(t)

	mockCalls := grpcMock.NewMockAuthenticationServiceClient(ctrl)

	g.client.authgRPClient = mockCalls

	req := services.RegisterUserRequest{
		FullName: "Jane",
		Email: "jane@gmail.com",
		Password: "password",
		PaydUsername: "username",
		PaydAccountID: "account_id",
		UsernameApiKey: "user_key",
		PasswordApiKey: "pass_key",
	}

	tests := []struct{
		name string
		buildStubs func(*grpcMock.MockAuthenticationServiceClient)
		want *pb.RegisterUserResponse
		wantStatusCode int
		wantErr bool
	}{
		{
			name: "success",
			buildStubs: func(mockCalls *grpcMock.MockAuthenticationServiceClient) {
				mockCalls.EXPECT().RegisterUser(gomock.Any(), gomock.Eq(&pb.RegisterUserRequest{
					Fullname:           req.FullName,
					Email:              req.Email,
					Password:           req.Password,
					PaydUsername:       req.PaydUsername,
					PaydPasswordApiKey: req.PasswordApiKey,
					PaydUsernameApiKey: req.UsernameApiKey,
					PaydAccountId: req.PaydAccountID,
				})).Return(&pb.RegisterUserResponse{
					Fullname: req.FullName,
					Email: req.Email,
					CreatedAt: timestamppb.New(TestTime),
				}, nil).Times(1)
			},
			want: &pb.RegisterUserResponse{
				Fullname: req.FullName,
				Email: req.Email,
				CreatedAt: timestamppb.New(TestTime),
			},
			wantStatusCode: http.StatusOK,
			wantErr: false,
		},
		{
			name: "internal error",
			buildStubs: func(mockCalls *grpcMock.MockAuthenticationServiceClient) {
				mockCalls.EXPECT().RegisterUser(gomock.Any(), gomock.Eq(&pb.RegisterUserRequest{
					Fullname:           req.FullName,
					Email:              req.Email,
					Password:           req.Password,
					PaydUsername:       req.PaydUsername,
					PaydPasswordApiKey: req.PasswordApiKey,
					PaydUsernameApiKey: req.UsernameApiKey,
					PaydAccountId: req.PaydAccountID,
				})).Return(nil, status.Errorf(codes.Internal, "internal error")).
				Times(1)
			},
			want: nil,
			wantStatusCode: http.StatusInternalServerError,
			wantErr: true,
		},
		{
			name: "some unknown error",
			buildStubs: func(mockCalls *grpcMock.MockAuthenticationServiceClient) {
				mockCalls.EXPECT().RegisterUser(gomock.Any(), gomock.Eq(&pb.RegisterUserRequest{
					Fullname:           req.FullName,
					Email:              req.Email,
					Password:           req.Password,
					PaydUsername:       req.PaydUsername,
					PaydPasswordApiKey: req.PasswordApiKey,
					PaydUsernameApiKey: req.UsernameApiKey,
					PaydAccountId: req.PaydAccountID,
				})).Return(nil, errors.New("some unknown error")).
				Times(1)
			},
			want: nil,
			wantStatusCode: http.StatusInternalServerError,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(mockCalls)
			statusCode, msg := g.client.RegisterUserViagRPC(req)
			require.Equal(t, statusCode, tc.wantStatusCode)
			if !tc.wantErr {
				if !proto.Equal(msg["response"].(*pb.RegisterUserResponse), tc.want) {
					t.Errorf("register user got = %v, want %v", msg["response"], tc.want)
				}
			}
		})
	}
}

func TestGrpcClient_LoginUserViagRPC(t *testing.T) {
	g := NewTestGrpcClient()

	ctrl := gomock.NewController(t)

	mockCalls := grpcMock.NewMockAuthenticationServiceClient(ctrl)

	g.client.authgRPClient = mockCalls

	req := services.LoginUserRequest{
		Email: "jane@gmail.com",
		Password: "secret",
	}

	rsp := &pb.LoginUserResponse{
		AccessToken: "token",
		ExpirationAt: timestamppb.New(TestTime),
		Data: &pb.RegisterUserResponse{
			Fullname: "jane doe",
			Email: req.Email,
			CreatedAt: timestamppb.New(TestTime),
		},
	}

	tests := []struct{
		name string
		buildStubs func(*grpcMock.MockAuthenticationServiceClient)
		want *pb.LoginUserResponse
		wantStatusCode int
		wantErr bool
	}{
		{
			name: "success",
			buildStubs: func(mockCalls *grpcMock.MockAuthenticationServiceClient) {
				mockCalls.EXPECT().LoginUser(gomock.Any(), gomock.Eq(&pb.LoginUserRequest{
					Email:              req.Email,
					Password:           req.Password,
				})).Return(rsp, nil).Times(1)
			},
			want: rsp,
			wantStatusCode: http.StatusOK,
			wantErr: false,
		},
		{
			name: "internal error",
			buildStubs: func(mockCalls *grpcMock.MockAuthenticationServiceClient) {
				mockCalls.EXPECT().LoginUser(gomock.Any(), gomock.Eq(&pb.LoginUserRequest{
					Email:              req.Email,
					Password:           req.Password,
				})).Return(nil, status.Errorf(codes.Internal, "internal error")).
				Times(1)
			},
			want: nil,
			wantStatusCode: http.StatusInternalServerError,
			wantErr: true,
		},
		{
			name: "some unknown error",
			buildStubs: func(mockCalls *grpcMock.MockAuthenticationServiceClient) {
				mockCalls.EXPECT().LoginUser(gomock.Any(), gomock.Eq(&pb.LoginUserRequest{
					Email:              req.Email,
					Password:           req.Password,
				})).Return(nil, errors.New("some unknown error")).
				Times(1)
			},
			want: nil,
			wantStatusCode: http.StatusInternalServerError,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(mockCalls)
			statusCode, msg := g.client.LoginUserViagRPC(req)
			require.Equal(t, statusCode, tc.wantStatusCode)
			if !tc.wantErr {
				if !proto.Equal(msg["response"].(*pb.LoginUserResponse), tc.want) {
					t.Errorf("login user got = %v, want %v", msg["response"], tc.want)
				}
			}
		})
	}
}