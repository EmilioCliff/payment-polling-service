package gRPC

import (
	"log"
	"net"

	"github.com/EmilioCliff/payment-polling-app/authentication-service/internal/postgres"
	"github.com/EmilioCliff/payment-polling-service/shared-grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GrpcServer struct {
	pb.UnimplementedAuthenticationServiceServer
	store *postgres.Store

	Addr string
}

func NewGrpcServer() (*GrpcServer, error) {
	store, err := postgres.NewStore()
	if err != nil {
		return nil, err
	}

	return &GrpcServer{
		store: store,
		Addr:  store.GRPC_ADDR,
	}, nil
}

func (s *GrpcServer) Start() {
	grpcServer := grpc.NewServer()

	pb.RegisterAuthenticationServiceServer(grpcServer, s)

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		log.Printf("Failed to start grpc server on port: %s", err)
		return
	}

	if err = grpcServer.Serve(listener); err != nil {
		log.Printf("Failed to start gRPC server: %s", err)
		return
	}
}
