package gRPC

import (
	"net/http"

	"github.com/EmilioCliff/payment-polling-app/gateway-service/internal/services"
	"github.com/EmilioCliff/payment-polling-service/shared-grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

var _ services.GrpcInterface = (*GrpcClient)(nil)

// type for dialing grpc server implemented to allow mocking.
type GrpcDialFunc func(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error)

type GrpcClient struct {
	authgRPClient pb.AuthenticationServiceClient
	dialFunc      GrpcDialFunc
}

func NewGrpcService() *GrpcClient {
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

func grpcCodeConvert(code codes.Code) int {
	switch code {
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusForbidden
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.PermissionDenied:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
