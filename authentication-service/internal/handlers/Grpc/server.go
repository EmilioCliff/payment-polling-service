package Grpc

import (
	"fmt"
	"net"
	"sync"

	"github.com/EmilioCliff/payment-polling-app/authentication-service/internal/repository"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/pkg"
	"github.com/EmilioCliff/payment-polling-service/shared-grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	pb.UnimplementedAuthenticationServiceServer
	gRPCServer *grpc.Server
	config     pkg.Config
	maker      pkg.JWTMaker
	// shutdownCh chan struct{}
	mu sync.Mutex

	UserRepository repository.UserRepository
}

func NewGRPCServer(config pkg.Config, tokenMaker pkg.JWTMaker) *GRPCServer {
	return &GRPCServer{
		config: config,
		maker:  tokenMaker,
		// shutdownCh: make(chan struct{}),
	}
}

func (s *GRPCServer) Start(port string) error {
	s.mu.Lock()
	s.gRPCServer = grpc.NewServer()
	s.mu.Unlock()

	pb.RegisterAuthenticationServiceServer(s.gRPCServer, s)

	reflection.Register(s.gRPCServer)

	listener, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("failed to start grpc server on port: %s", err)
	}

	return s.gRPCServer.Serve(listener)
}

func (s *GRPCServer) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.gRPCServer != nil {
		s.gRPCServer.GracefulStop()
	}
}
