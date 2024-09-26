package gRPC

import (
	"context"
	"net/http"
	"time"

	"github.com/EmilioCliff/payment-polling-app/gateway-service/internal/services"
	"github.com/EmilioCliff/payment-polling-service/shared-grpc/pb"
	"google.golang.org/grpc/status"
)

func (g *GrpcClient) RegisterUserViagRPC(req services.RegisterUserRequest) (int, services.RegisterUserResponse) {
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	rsp, err := g.authgRPClient.RegisterUser(c, &pb.RegisterUserRequest{
		Fullname:           req.FullName,
		Email:              req.Email,
		Password:           req.Password,
		PaydUsername:       req.PaydUsername,
		PaydPasswordApiKey: req.PasswordApiKey,
		PaydUsernameApiKey: req.UsernameApiKey,
		PaydAccountId:      req.PaydAccountID,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			code := grpcCodeConvert(st.Code())
			grpcMessage := st.Message()

			return code, services.RegisterUserResponse{Message: grpcMessage, StatusCode: code}
		}

		return http.StatusInternalServerError, services.RegisterUserResponse{Message: "internal error", StatusCode: http.StatusInternalServerError}
	}

	return http.StatusOK, services.RegisterUserResponse{
		FullName:  rsp.GetFullname(),
		Email:     rsp.GetEmail(),
		CreatedAt: rsp.GetCreatedAt().AsTime(),
	}
}

func (g *GrpcClient) LoginUserViagRPC(req services.LoginUserRequest) (int, services.LoginUserResponse) {
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	rsp, err := g.authgRPClient.LoginUser(c, &pb.LoginUserRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			code := grpcCodeConvert(st.Code())
			grpcMessage := st.Message()

			return code, services.LoginUserResponse{Message: grpcMessage, StatusCode: code}
		}

		return http.StatusInternalServerError, services.LoginUserResponse{Message: "internal error", StatusCode: http.StatusInternalServerError}
	}

	return http.StatusOK, services.LoginUserResponse{
		AccessToken:  rsp.GetAccessToken(),
		FullName:     rsp.GetData().GetFullname(),
		Email:        rsp.GetData().GetEmail(),
		ExpirationAt: rsp.GetExpirationAt().AsTime(),
		CreatedAt:    rsp.GetData().GetCreatedAt().AsTime(),
	}
}
