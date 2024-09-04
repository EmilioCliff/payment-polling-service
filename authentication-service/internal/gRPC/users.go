package gRPC

import (
	"context"
	"database/sql"
	"fmt"

	pb "github.com/EmilioCliff/payment-polling-service/shared-grpc/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *GrpcServer) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	reqData := s.convertGrpcRegisterReq(req)

	rspData, err := s.store.RegisterUser(ctx, reqData)
	if err != nil {
		grpcCode := s.convertPkgError(err)
		return nil, status.Errorf(grpcCode, fmt.Sprintf("error on register user: %s", err.Message))
	}

	rsp := s.convertRegisterResp(rspData)
	return rsp, nil
}

func (s *GrpcServer) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	reqData := s.convertGrpcLoginReq(req)

	rspData, err := s.store.LoginUser(ctx, reqData)
	if err != nil {
		grpcCode := s.convertPkgError(err)
		return nil, status.Errorf(grpcCode, fmt.Sprintf("error on logging user: %s", err.Message))
	}

	rsp := s.convertLoginResp(rspData)

	return rsp, nil
}

func (s *GrpcServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := s.store.Queries.GetUser(ctx, req.GetUserId())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	rsp := s.convertGetUserResp(user)

	return rsp, nil
}
