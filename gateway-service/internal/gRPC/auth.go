package gRPC

import (
	"context"
	"net/http"
	"time"

	"github.com/EmilioCliff/payment-polling-app/gateway-service/internal/services"
	"github.com/EmilioCliff/payment-polling-app/gateway-service/pkg"
	"github.com/EmilioCliff/payment-polling-service/shared-grpc/pb"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"
)

func (g *GrpcClient) RegisterUserViagRPC(req services.RegisterUserRequest) (int, gin.H) {
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	rsp, err := g.authgRPClient.RegisterUser(c, &pb.RegisterUserRequest{
		Fullname:           req.FullName,
		Email:              req.Email,
		Password:           req.Password,
		PaydUsername:       req.PaydUsername,
		PaydPasswordApiKey: req.PasswordApiKey,
		PaydUsernameApiKey: req.UsernameApiKey,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			grpcCode := st.Code()
			grpcMessage := st.Message()
			return g.grpcErrorResponse(grpcCode, grpcMessage)
		} else {
			return http.StatusInternalServerError, pkg.ErrorResponse(err, "Non-gRPC error occurred")
		}
	}

	return http.StatusOK, gin.H{"response": rsp}
}

func (g *GrpcClient) LoginUserViagRPC(req services.LoginUserRequest) (int, gin.H) {
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	rsp, err := g.authgRPClient.LoginUser(c, &pb.LoginUserRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			grpcCode := st.Code()
			grpcMessage := st.Message()
			return g.grpcErrorResponse(grpcCode, grpcMessage)
		} else {
			return http.StatusInternalServerError, pkg.ErrorResponse(err, "Non-gRPC error occurred")
		}
	}

	return http.StatusOK, gin.H{"response": rsp}
}
