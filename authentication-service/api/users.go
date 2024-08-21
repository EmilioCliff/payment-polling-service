package api

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	db "github.com/EmilioCliff/payment-polling-app/authentication-service/db/sqlc"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/utils"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
)

type registerResponse struct {
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (server *Server) registerUserGeneral(req registerUserRequest, ctx context.Context) (registerResponse, generalErrorResponse) {
	hashPassword, err := utils.GenerateHashPassword(req.Password, server.config.HASH_COST)
	if err != nil {
		return registerResponse{}, errorHelper("error hashing password", err, http.StatusInternalServerError, codes.Internal)
	}

	user, err := server.store.RegisterUser(ctx, db.RegisterUserParams{
		FullName: req.FullName,
		Email:    req.Email,
		Password: hashPassword,
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return registerResponse{}, errorHelper("user already exists", err, http.StatusForbidden, codes.AlreadyExists)
			}
		}

		return registerResponse{}, errorHelper("error creating user", err, http.StatusInternalServerError, codes.Internal)
	}

	return registerResponse{
		FullName:  user.FullName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, generalErrorResponse{Status: true}
}

type loginResponse struct {
	AccessToken  string    `json:"access_token"`
	ExpirationAt time.Time `json:"expiration_at"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"created_at"`
}

func (server *Server) loginUserGeneral(req loginUserRequest, ctx context.Context) (loginResponse, generalErrorResponse) {
	user, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return loginResponse{}, errorHelper("user not found", err, http.StatusNotFound, codes.NotFound)
		}
		return loginResponse{}, errorHelper("error getting user by email", err, http.StatusInternalServerError, codes.Internal)
	}

	err = utils.ComparePasswordAndHash(user.Password, req.Password)
	if err != nil {
		return loginResponse{}, errorHelper("invalid credentials", err, http.StatusUnauthorized, codes.Unauthenticated)
	}

	accessToken, err := server.maker.CreateToken(user.Email, server.config.TOKEN_DURATION)
	if err != nil {
		return loginResponse{}, errorHelper("error creating access token", err, http.StatusInternalServerError, codes.Internal)
	}

	return loginResponse{
		AccessToken:  accessToken,
		ExpirationAt: time.Now().Add(server.config.TOKEN_DURATION),
		FullName:     user.FullName,
		Email:        user.Email,
		CreatedAt:    user.CreatedAt,
	}, generalErrorResponse{Status: true}
}
