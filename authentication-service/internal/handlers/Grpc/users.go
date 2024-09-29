package Grpc

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/EmilioCliff/payment-polling-app/authentication-service/internal/repository"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/pkg"
	"github.com/EmilioCliff/payment-polling-service/shared-grpc/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *GRPCServer) RegisterUser(
	ctx context.Context,
	req *pb.RegisterUserRequest,
) (*pb.RegisterUserResponse, error) {
	user, err := s.UserRepository.CreateUser(ctx, repository.User{
		FullName:        req.GetFullname(),
		Email:           req.GetEmail(),
		Password:        req.GetPassword(),
		PaydUsername:    req.GetPaydUsername(),
		PaydUsernameKey: req.GetPaydUsernameApiKey(),
		PaydPasswordKey: req.GetPaydPasswordApiKey(),
		PaydAccountID:   req.GetPaydAccountId(),
	})
	if err != nil {
		grpcCode := convertPkgError(pkg.ErrorCode(err))

		return nil, status.Errorf(
			grpcCode,
			"%v",
			fmt.Sprintf("error on register user: %v", pkg.ErrorMessage(err)),
		)
	}

	return &pb.RegisterUserResponse{
		Fullname:  user.FullName,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}, nil
}

func (s *GRPCServer) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := s.UserRepository.GetUser(ctx, req.GetEmail())
	if err != nil {
		grpcCode := convertPkgError(pkg.ErrorCode(err))

		return nil, status.Errorf(
			grpcCode,
			"%v",
			fmt.Sprintf("error on logging user: %v", pkg.ErrorMessage(err)),
		)
	}

	err = pkg.ComparePasswordAndHash(user.Password, req.GetPassword())
	if err != nil {
		grpcCode := convertPkgError(pkg.ErrorCode(err))

		return nil, status.Errorf(
			grpcCode,
			"%v",
			fmt.Sprintf("Error comparing passwords: %v", pkg.ErrorMessage(err)),
		)
	}

	accessToken, err := s.maker.CreateToken(user.Email, user.ID, s.config.TOKEN_DURATION)
	if err != nil {
		grpcCode := convertPkgError(pkg.ErrorCode(err))

		return nil, status.Errorf(
			grpcCode,
			"%v",
			fmt.Sprintf("Error creating token: %v", pkg.ErrorMessage(err)),
		)
	}

	return &pb.LoginUserResponse{
		AccessToken:  accessToken,
		ExpirationAt: timestamppb.New(time.Now().Add(s.config.TOKEN_DURATION)),
		Data: &pb.RegisterUserResponse{
			Fullname:  user.FullName,
			Email:     user.Email,
			CreatedAt: timestamppb.New(user.CreatedAt),
		},
	}, nil
}

func (s *GRPCServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := s.UserRepository.GetUser(ctx, req.GetEmail())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}

		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	return &pb.GetUserResponse{
		UserId:          user.ID,
		PaydUsername:    user.PaydUsername,
		PaydUsernameKey: user.PaydUsernameKey,
		PaydPasswordKey: user.PaydPasswordKey,
		PaydAccountId:   user.PaydAccountID,
	}, nil
}

func convertPkgError(err string) codes.Code {
	switch err {
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
