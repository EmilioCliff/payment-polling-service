package gRPC

import (
	"net/http"

	"github.com/EmilioCliff/payment-polling-app/gateway-service/internal/services"
	"github.com/EmilioCliff/payment-polling-service/shared-grpc/pb"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

var _ services.GrpcInterface = (*GrpcClient)(nil)

type GrpcClient struct {
	authgRPClient pb.AuthenticationServiceClient
}

func NewGrpcClient() *GrpcClient {
	return &GrpcClient{}
}

func (g *GrpcClient) Start(grpcPort string) error {
	gRPCconn, err := grpc.NewClient(grpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	g.authgRPClient = pb.NewAuthenticationServiceClient(gRPCconn)
	return nil
}

func (g *GrpcClient) grpcErrorResponse(code codes.Code, msg string) (int, gin.H) {
	switch code {
	case codes.Internal:
		return http.StatusInternalServerError, gin.H{
			"message": msg,
		}
	case codes.NotFound:
		return http.StatusNotFound, gin.H{
			"message": msg,
		}
	case codes.AlreadyExists:
		return http.StatusForbidden, gin.H{
			"message": msg,
		}
	case codes.Unauthenticated:
		return http.StatusUnauthorized, gin.H{
			"message": msg,
		}
	case codes.PermissionDenied:
		return http.StatusForbidden, gin.H{
			"message": msg,
		}
	default:
		return http.StatusInternalServerError, gin.H{
			"message": msg,
		}
	}
}
