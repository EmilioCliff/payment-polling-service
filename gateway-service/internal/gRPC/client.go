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

type GrpcDialFunc func(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error)

type GrpcClient struct {
	authgRPClient pb.AuthenticationServiceClient
	dialFunc      GrpcDialFunc
}

func NewGrpcClient() *GrpcClient {
	return &GrpcClient{
		dialFunc: grpc.NewClient,
	}
}

func (g *GrpcClient) Start(grpcPort string) error {
	gRPCconn, err := g.dialFunc(grpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
