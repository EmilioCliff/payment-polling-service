package gRPC

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type TestGrpcClient struct {
	client *GrpcClient
}

func NewTestGrpcClient() *TestGrpcClient {

	g := &TestGrpcClient{
		client: NewGrpcClient(),
	}

	g.client.dialFunc = func(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
		if target == "INVALID_PORT" {return nil, errors.New("invalid port")}
		return &grpc.ClientConn{}, nil
	}

	return g
}


func TestGrpcClient_Start(t *testing.T) {

	g := NewTestGrpcClient()

	tests := []struct{
		name string
		port string
		wantErr bool
	}{
		{
			name: "valid port",
			port: "0.0.0.0:5050",
			wantErr: false,
		},
		{
			name: "invalid port",
			port: "INVALID_PORT",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := g.client.Start(tc.port)
			if(err != nil) != tc.wantErr  {
				t.Errorf("Start() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestGrpcClient_grpcErrorResponse(t *testing.T) {
	g := NewTestGrpcClient()

	tests := []struct{
		name string
		code codes.Code
		msg string
		wantStatusCode int
		wantMessage string
	}{
		{
			name: "internal error",
			code: codes.Internal,
			msg: "internal error",
			wantStatusCode: http.StatusInternalServerError,
			wantMessage: "internal error",
		},
		{
			name: "not found",
			code: codes.NotFound,
			msg: "not found",
			wantStatusCode: http.StatusNotFound,
			wantMessage: "not found",
		},
		{
			name: "already exists",
			code: codes.AlreadyExists,
			msg: "already exists",
			wantStatusCode: http.StatusForbidden,
			wantMessage: "already exists",
		},
		{
			name: "unathorized",
			code: codes.Unauthenticated,
			msg: "unathorized",
			wantStatusCode: http.StatusUnauthorized,
			wantMessage: "unathorized",
		},
		{
			name: "permission denied",
			code: codes.PermissionDenied,
			msg: "permission denied",
			wantStatusCode: http.StatusForbidden,
			wantMessage: "permission denied",
		},
		{
			name: "system error",
			code: codes.ResourceExhausted, // any code
			msg: "system error",
			wantStatusCode: http.StatusInternalServerError,
			wantMessage: "system error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			statusCode, msg := g.client.grpcErrorResponse(tc.code, tc.msg)
			require.Equal(t, statusCode, tc.wantStatusCode)
			require.Equal(t, msg, gin.H{"message": tc.msg})
		})
	}
}