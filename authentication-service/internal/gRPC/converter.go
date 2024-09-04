package gRPC

import (
	"github.com/EmilioCliff/payment-polling-app/authentication-service/internal/postgres"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/internal/postgres/generated"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/pkg"
	"github.com/EmilioCliff/payment-polling-service/shared-grpc/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *GrpcServer) convertGrpcRegisterReq(req *pb.RegisterUserRequest) postgres.RegisterUserRequest {
	return postgres.RegisterUserRequest{
		FullName:       req.GetFullname(),
		PaydUsername:   req.GetPaydUsername(),
		Email:          req.GetEmail(),
		Password:       req.GetPassword(),
		UsernameApiKey: req.GetPaydUsernameApiKey(),
		PasswordApiKey: req.GetPaydPasswordApiKey(),
	}
}

func (s *GrpcServer) convertRegisterResp(rsp *postgres.RegisterUserResponse) *pb.RegisterUserResponse {
	return &pb.RegisterUserResponse{
		Fullname:  rsp.FullName,
		Email:     rsp.Email,
		CreatedAt: timestamppb.New(rsp.CreatedAt),
	}
}

func (s *GrpcServer) convertGrpcLoginReq(req *pb.LoginUserRequest) postgres.LoginUserRequest {
	return postgres.LoginUserRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}
}

func (s *GrpcServer) convertLoginResp(rsp *postgres.LoginUserResponse) *pb.LoginUserResponse {
	return &pb.LoginUserResponse{
		AccessToken:  rsp.AccessToken,
		ExpirationAt: timestamppb.New(rsp.ExpirationAt),
		Data: s.convertRegisterResp(&postgres.RegisterUserResponse{
			FullName:  rsp.FullName,
			Email:     rsp.Email,
			CreatedAt: rsp.CreatedAt,
		}),
	}
}

func (s *GrpcServer) convertGetUserResp(rsp generated.User) *pb.GetUserResponse {
	return &pb.GetUserResponse{
		PaydUsername:    rsp.PaydUsername,
		PaydUsernameKey: rsp.PaydUsernameKey,
		PaydPasswordKey: rsp.PaydPasswordKey,
	}
}

func (s *GrpcServer) convertPkgError(err *pkg.Error) codes.Code {
	switch err.Code {
	case pkg.ALREADY_EXISTS_ERROR:
		return codes.AlreadyExists
	case pkg.INTERNAL_ERROR:
		return codes.Internal
	case pkg.INVALID_ERROR:
		return codes.InvalidArgument
	case pkg.NOT_FOUND_ERROR:
		return codes.NotFound
	case pkg.NOT_IMPLEMENTED_ERROR:
		return codes.Unimplemented
	case pkg.AUTHENTICATION_ERROR:
		return codes.Unauthenticated
	default:
		return codes.Internal
	}
}
