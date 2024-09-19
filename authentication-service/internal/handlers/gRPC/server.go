package gRPC

import (
	"fmt"
	"net"

	"github.com/EmilioCliff/payment-polling-app/authentication-service/internal/repository"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/pkg"
	pb "github.com/EmilioCliff/payment-polling-service/shared-grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	gRPCServer *grpc.Server
	pb.UnimplementedAuthenticationServiceServer
	config pkg.Config
	maker pkg.JWTMaker

	UserRepository repository.UserRepository
}

func NewGRPCServer(config pkg.Config, tokenMaker pkg.JWTMaker) *GRPCServer {
	return &GRPCServer{
		config: config,
		maker: tokenMaker,
	}
}

func (s *GRPCServer) Start(port string) error {
	s.gRPCServer = grpc.NewServer()

	pb.RegisterAuthenticationServiceServer(s.gRPCServer, s)

	reflection.Register(s.gRPCServer)

	listener, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("Failed to start grpc server on port: %s", err)
	}

	if err = s.gRPCServer.Serve(listener); err != nil {
		return fmt.Errorf("Failed to start gRPC server: %s", err)
	}

	return nil
}

func (s *GRPCServer) Stop() {
	s.gRPCServer.GracefulStop()
}