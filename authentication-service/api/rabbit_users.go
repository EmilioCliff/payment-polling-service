package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	db "github.com/EmilioCliff/payment-polling-app/authentication-service/db/sqlc"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/utils"
	"github.com/lib/pq"
)

func (server *Server) rabbitRegisterUser(data map[string]interface{}) []byte {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	hashPassword, err := utils.GenerateHashPassword(data["password"].(string), server.config.HASH_COST)
	if err != nil {
		return []byte("error hashing password")
	}

	user, err := server.store.RegisterUser(ctx, db.RegisterUserParams{
		FullName: data["full_name"].(string),
		Email:    data["email"].(string),
		Password: hashPassword,
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return []byte("user already exists")
			}
		}
		return []byte("failed to create user")
	}

	responseData := registerUserResponse{
		Status: true,
		Data: registerResponse{
			FullName:  user.FullName,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	}

	rsp, err := json.Marshal(responseData)
	if err != nil {
		return []byte("error marshallinggg response")
	}

	return rsp
}

func (server *Server) rabbitLoginUser(data map[string]interface{}) []byte {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	user, err := server.store.GetUserByEmail(ctx, data["email"].(string))
	if err != nil {
		if err == sql.ErrNoRows {
			return []byte("user not found")
		}
		return []byte("error getting user by email")
	}

	err = utils.ComparePasswordAndHash(user.Password, data["password"].(string))
	if err != nil {
		return []byte("invalid credentials")
	}

	accessToken, err := server.maker.CreateToken(user.Email, server.config.TOKEN_DURATION)
	if err != nil {
		return []byte("error creating access token")
	}

	response := loginUserResponse{
		Status: true,
		Data: loginResponse{
			AccessToken:  accessToken,
			ExpirationAt: time.Now().Add(server.config.TOKEN_DURATION),
			FullName:     user.FullName,
			Email:        user.Email,
			CreatedAt:    user.CreatedAt,
		},
	}

	rsp, err := json.Marshal(response)
	if err != nil {
		return []byte("error marshalling response")
	}

	return rsp
}
