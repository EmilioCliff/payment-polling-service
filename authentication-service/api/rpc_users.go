package api

import (
	"context"
	"database/sql"
	"log"
	"time"

	db "github.com/EmilioCliff/payment-polling-app/authentication-service/db/sqlc"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/pb"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/utils"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	hashPassword, err := utils.GenerateHashPassword(req.GetPassword(), server.config.HASH_COST)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}

	user, err := server.store.RegisterUser(ctx, db.RegisterUserParams{
		FullName: req.GetFullname(),
		Email:    req.GetEmail(),
		Password: hashPassword,
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			log.Println(pqErr.Code.Name())
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "user already exists: %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	rsp := &pb.RegisterUserResponse{
		Fullname:  user.FullName,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}

	return rsp, nil
}
func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := server.store.GetUserByEmail(ctx, req.GetEmail())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user does not exist: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to get user: %s", err)
	}

	err = utils.ComparePasswordAndHash(user.Password, req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "credential do not match: %s", err)
	}

	accessToken, err := server.maker.CreateToken(user.Email, server.config.TOKEN_DURATION)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error generating token: %s", err)
	}

	rsp := &pb.LoginUserResponse{
		AccessToken:  accessToken,
		ExpirationAt: timestamppb.New(time.Now().Add(server.config.TOKEN_DURATION)),
		Data: &pb.RegisterUserResponse{
			Fullname:  user.FullName,
			Email:     user.Email,
			CreatedAt: timestamppb.New(user.CreatedAt),
		},
	}

	return rsp, nil
}
