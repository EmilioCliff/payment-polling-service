package gRPC

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type TestGrpcClient struct {
	client *GrpcClient
}

func NewTestGrpcClient() *TestGrpcClient {
	g := &TestGrpcClient{
		client: NewGrpcService(),
	}

	// mock the dial function
	g.client.dialFunc = func(target string, _ ...grpc.DialOption) (*grpc.ClientConn, error) {
		if target == "INVALID_PORT" {
			return nil, errors.New("invalid port")
		}

		return &grpc.ClientConn{}, nil
	}

	return g
}

func TestGrpcClient_Start(t *testing.T) {
	g := NewTestGrpcClient()

	tests := []struct {
		name    string
		port    string
		wantErr bool
	}{
		{
			name:    "valid port",
			port:    "0.0.0.0:5050",
			wantErr: false,
		},
		{
			name:    "invalid port",
			port:    "INVALID_PORT",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := g.client.Start(tc.port)
			if (err != nil) != tc.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestGrpcClient_grpcCodeConvert(t *testing.T) {
	tests := []struct {
		name       string
		code       codes.Code
		statusCode int
	}{
		{
			name:       "internal error",
			code:       codes.Internal,
			statusCode: http.StatusInternalServerError,
		},
		{
			name:       "not found",
			code:       codes.NotFound,
			statusCode: http.StatusNotFound,
		},
		{
			name:       "already exists",
			code:       codes.AlreadyExists,
			statusCode: http.StatusForbidden,
		},
		{
			name:       "unathorized",
			code:       codes.Unauthenticated,
			statusCode: http.StatusUnauthorized,
		},
		{
			name:       "permission denied",
			code:       codes.PermissionDenied,
			statusCode: http.StatusForbidden,
		},
		{
			name:       "system error",
			code:       codes.ResourceExhausted, // any code
			statusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			statusCode := grpcCodeConvert(tc.code)
			require.Equal(t, statusCode, tc.statusCode)
		})
	}
}
