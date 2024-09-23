package Grpc

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/EmilioCliff/payment-polling-app/authentication-service/internal/mock"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/pkg"
)

type TestGRPCServer struct {
	server         *GRPCServer
	UserRepository mock.MockUsersRepositry
}

func NewTestGRPCServer() *TestGRPCServer {
	bits := 2048

	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		panic(err)
	}

	publicKey := &privateKey.PublicKey

	s := &TestGRPCServer{
		server: NewGRPCServer(pkg.Config{
			TOKEN_DURATION: time.Second,
		}, pkg.JWTMaker{
			PublicKey:  publicKey,
			PrivateKey: privateKey,
		}),
	}

	s.server.UserRepository = &s.UserRepository

	return s
}

func TestGRPCServer_Start(t *testing.T) {
	tests := []struct {
		name    string
		port    string
		wantErr bool
	}{
		{
			name:    "Valid Port",
			port:    "0.0.0.0:5050",
			wantErr: false,
		},
		{
			name:    "Invalid Port",
			port:    "INVALID_PORT",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := NewTestGRPCServer()
			errCh := make(chan error)

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			go func(port string) {
				errCh <- s.server.Start(port)
			}(tc.port)

			select {
			case err := <-errCh:
				if (err != nil) != tc.wantErr {
					t.Errorf("GRPCServer.Run() error = %v, wantErr %v", err, tc.wantErr)
				}
			case <-ctx.Done():
				if tc.wantErr {
					t.Errorf("GRPCServer.Run() expected error for port %s, but got none", tc.port)
				}
			}

			s.server.Stop()
		})
	}
}
