package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"testing"

	"github.com/EmilioCliff/payment-polling-app/payment-service/internal/repository"
	"github.com/EmilioCliff/payment-polling-app/payment-service/internal/services"
	"github.com/EmilioCliff/payment-polling-app/payment-service/pkg"
	"github.com/EmilioCliff/payment-polling-service/shared-grpc/mockpb"
	"github.com/EmilioCliff/payment-polling-service/shared-grpc/pb"
	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
)

func mockDistributeSendPaymentRequestTaskFunc(ctx context.Context, payload services.SendPaymentWithdrawalRequestPayload, opt ...asynq.Option) error {
	if payload.Amount == 32 {
		return errors.New("invalid payload")
	}

	return nil
}

func mockDistributeSendWithdrawalRequestTaskFunc(
	ctx context.Context,
	payload services.SendPaymentWithdrawalRequestPayload,
	opt ...asynq.Option,
) error {
	if payload.Amount == 32 {
		log.Println("here")

		return errors.New("invalid payload")
	}

	return nil
}

func TestRabbitConn_handleInitiatePayment(t *testing.T) {
	r := NewTestRabbitHandler()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedClient := mockpb.NewMockAuthenticationServiceClient(ctrl)

	r.rabbit.client = mockedClient
	r.TastDistributor.DistributeSendPaymentRequestTaskFunc = mockDistributeSendPaymentRequestTaskFunc
	r.TastDistributor.DistributeSendWithdrawalRequestTaskFunc = mockDistributeSendWithdrawalRequestTaskFunc

	pbGetUserStub := func(ctx context.Context, in *pb.GetUserRequest, opts ...grpc.CallOption) (*pb.GetUserResponse, error) {
		emailEncrypted, err := pkg.Encrypt(in.GetEmail(), []byte(r.rabbit.config.ENCRYPTION_KEY))
		require.NoError(t, err)

		return &pb.GetUserResponse{
			UserId:          32,
			PaydUsername:    "test",
			PaydUsernameKey: emailEncrypted,
			PaydPasswordKey: emailEncrypted,
			PaydAccountId:   "test",
		}, nil
	}

	tests := []struct {
		name         string
		req          initiatePaymentRequest
		buildPbStubs func(*mockpb.MockAuthenticationServiceClient, string, any)
		wantRsp      any
		wantErr      bool
	}{
		{
			name: "payment success",
			req: initiatePaymentRequest{
				Email:       "test",
				Action:      "payment",
				Amount:      100,
				PhoneNumber: "test",
				NetworkCode: "test",
				Naration:    "test",
			},
			buildPbStubs: func(mockedClient *mockpb.MockAuthenticationServiceClient, email string, pbGetUserStub any) {
				mockedClient.EXPECT().GetUser(gomock.Any(), &pb.GetUserRequest{Email: email}).
					DoAndReturn(pbGetUserStub).Times(1)
			},
			wantRsp: initiatePaymentResponse{
				TransactionID: "test",
				PaymentStatus: false,
				Action:        "payment",
			},
			wantErr: false,
		},
		{
			name: "withdrawal success",
			req: initiatePaymentRequest{
				Email:       "test",
				Action:      "withdrawal",
				Amount:      100,
				PhoneNumber: "test",
				NetworkCode: "test",
				Naration:    "test",
			},
			buildPbStubs: func(mockedClient *mockpb.MockAuthenticationServiceClient, email string, pbGetUserStub any) {
				mockedClient.EXPECT().GetUser(gomock.Any(), &pb.GetUserRequest{Email: email}).
					DoAndReturn(pbGetUserStub).Times(1)
			},
			wantRsp: initiatePaymentResponse{
				TransactionID: "test",
				PaymentStatus: false,
				Action:        "withdrawal",
			},
			wantErr: false,
		},
		{
			name: "invalid action",
			req: initiatePaymentRequest{
				Email:       "test",
				Action:      "invalid",
				Amount:      100,
				PhoneNumber: "test",
				NetworkCode: "test",
				Naration:    "test",
			},
			buildPbStubs: func(mockedClient *mockpb.MockAuthenticationServiceClient, email string, pbGetUserStub any) {
				mockedClient.EXPECT().GetUser(gomock.Any(), &pb.GetUserRequest{Email: email}).
					DoAndReturn(pbGetUserStub).Times(1)
			},
			wantRsp: errorResponse{
				Status:  http.StatusBadRequest,
				Message: "invalid action: invalid",
			},
			wantErr: true,
		},
		{
			name: "no user found",
			req: initiatePaymentRequest{
				Email:       "test",
				Action:      "withdrawal",
				Amount:      100,
				PhoneNumber: "test",
				NetworkCode: "test",
				Naration:    "test",
			},
			buildPbStubs: func(mockedClient *mockpb.MockAuthenticationServiceClient, email string, _ any) {
				mockedClient.EXPECT().GetUser(gomock.Any(), &pb.GetUserRequest{Email: email}).
					Return(nil, errors.New("user not found")).Times(1)
			},
			wantRsp: errorResponse{
				Status:  http.StatusInternalServerError,
				Message: "failed to get user data from auth: user not found",
			},
			wantErr: true,
		},
		{
			name: "task distribution error",
			req: initiatePaymentRequest{
				Email:       "test",
				Action:      "payment",
				Amount:      32,
				PhoneNumber: "test",
				NetworkCode: "test",
				Naration:    "test",
			},
			buildPbStubs: func(mockedClient *mockpb.MockAuthenticationServiceClient, email string, pbGetUserStub any) {
				mockedClient.EXPECT().GetUser(gomock.Any(), &pb.GetUserRequest{Email: email}).
					DoAndReturn(pbGetUserStub).Times(1)
			},
			wantRsp: errorResponse{
				Status:  http.StatusInternalServerError,
				Message: "failed to distribute payment task: invalid payload",
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildPbStubs(mockedClient, tc.req.Email, pbGetUserStub)

			rspBytes := r.rabbit.handleInitiatePayment(tc.req)

			if tc.wantErr {
				var rsp errorResponse

				err := json.Unmarshal(rspBytes, &rsp)
				require.NoError(t, err)

				rabbitError, ok := tc.wantRsp.(errorResponse)
				require.True(t, ok)

				require.Equal(t, rabbitError.Message, rsp.Message)
				require.Equal(t, rabbitError.Status, rsp.Status)
			} else {
				var rsp initiatePaymentResponse

				err := json.Unmarshal(rspBytes, &rsp)
				require.NoError(t, err)

				rabbitRsp, _ := tc.wantRsp.(initiatePaymentResponse)

				require.Equal(t, rabbitRsp.Action, rsp.Action)
			}
		})
	}
}

func mockPollingTransactionFunc(_ context.Context, id uuid.UUID) (*repository.Transaction, error) {
	if id == uuid.Nil {
		return nil, pkg.Errorf(pkg.INVALID_ERROR, "id cannot be empty")
	}

	return &repository.Transaction{
		TransactionID:      id,
		PaydTransactionRef: gofakeit.UUID(),
		UserID:             1,
		Message:            gofakeit.Sentence(10),
		Action:             "withdrawal",
		Amount:             100,
		PhoneNumber:        gofakeit.Phone(),
		NetworkCode:        "63902",
		Narration:          gofakeit.Sentence(10),
		Status:             false,
	}, nil
}

func TestRabbitConn_handlePollingTransaction(t *testing.T) {
	r := NewTestRabbitHandler()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r.TransactionRepository.PollingTransactionFunc = mockPollingTransactionFunc

	tests := []struct {
		name    string
		req     pollingTransactionRequest
		wantRsp any
		wantErr bool
	}{
		{
			name: "success",
			req: pollingTransactionRequest{
				TransactionId: gofakeit.UUID(),
				UserID:        1,
			},
			wantRsp: pollingTransactionResponse{
				Action:      "withdrawal",
				Amount:      100,
				NetworkCode: "63902",
			},
			wantErr: false,
		},
		{
			name: "invalid uuid",
			req: pollingTransactionRequest{
				TransactionId: "invalid uuid",
				UserID:        1,
			},
			wantRsp: errorResponse{
				Status:  http.StatusBadRequest,
				Message: "invalid transaction id: invalid UUID length: 12",
			},
			wantErr: true,
		},
		{
			name: "another user transaction",
			req: pollingTransactionRequest{
				TransactionId: gofakeit.UUID(),
				UserID:        32,
			},
			wantRsp: errorResponse{
				Status:  http.StatusUnauthorized,
				Message: "cannot access this transaction",
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rspBytes := r.rabbit.handlePollingTransaction(tc.req)

			if tc.wantErr {
				var rsp errorResponse

				err := json.Unmarshal(rspBytes, &rsp)
				require.NoError(t, err)

				rabbitRsp, ok := tc.wantRsp.(errorResponse)
				require.True(t, ok)

				require.Equal(t, rabbitRsp.Message, rsp.Message)
				require.Equal(t, rabbitRsp.Status, rsp.Status)
			} else {
				var rsp pollingTransactionResponse
				err := json.Unmarshal(rspBytes, &rsp)
				require.NoError(t, err)

				rabbitRsp, ok := tc.wantRsp.(pollingTransactionResponse)
				require.True(t, ok)

				require.Equal(t, rabbitRsp.Action, rsp.Action)
				require.Equal(t, rabbitRsp.Amount, rsp.Amount)
				require.Equal(t, rabbitRsp.NetworkCode, rsp.NetworkCode)
			}
		})
	}
}
