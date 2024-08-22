package api

import (
	"context"
	"fmt"

	pb "github.com/EmilioCliff/payment-polling-service/shared-grpc/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	reqData := registerUserRequest{
		FullName:       req.GetFullname(),
		Email:          req.GetEmail(),
		Password:       req.GetPassword(),
		PasswordApiKey: req.GetPaydPasswordApiKey(),
		UsernameApiKey: req.GetPaydUsernameApiKey(),
		PaydUsername:   req.GetPaydUsername(),
	}

	rspData, errRsp := server.registerUserGeneral(reqData, ctx)
	if !errRsp.Status {
		return nil, status.Errorf(errRsp.GrpcCode, fmt.Sprintf("%s: %s", errRsp.Error, errRsp.Message))
	}

	rsp := &pb.RegisterUserResponse{
		Fullname:  rspData.FullName,
		Email:     rspData.Email,
		CreatedAt: timestamppb.New(rspData.CreatedAt),
	}

	return rsp, nil
}
func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	reqData := loginUserRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	rspData, errRsp := server.loginUserGeneral(reqData, ctx)
	if !errRsp.Status {
		return nil, status.Errorf(errRsp.GrpcCode, fmt.Sprintf("%s: %s", errRsp.Error, errRsp.Message))
	}

	rsp := &pb.LoginUserResponse{
		AccessToken:  rspData.AccessToken,
		ExpirationAt: timestamppb.New(rspData.ExpirationAt),
		Data: &pb.RegisterUserResponse{
			Fullname:  rspData.FullName,
			Email:     rspData.Email,
			CreatedAt: timestamppb.New(rspData.CreatedAt),
		},
	}

	return rsp, nil
}

func (server *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserRequest, error) {

	server.store.GetUser(ctx, req.GetUserId())

	return nil, status.Errorf(codes.Unimplemented, "method GetUser not implemented")
}
