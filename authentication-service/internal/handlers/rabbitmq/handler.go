package rabbitmq

import (
	"context"
	"encoding/json"
	"time"

	"github.com/EmilioCliff/payment-polling-app/authentication-service/internal/repository"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/pkg"
)

type RegisterUserRequest struct {
	FullName        string `binding:"required" json:"full_name"`
	Email           string `binding:"required" json:"email"`
	Password        string `binding:"required" json:"password"`
	PaydUsername    string `binding:"required" json:"payd_username"`
	PaydAccountID   string `binding:"required" json:"payd_account_id"`
	PaydUsernameKey string `binding:"required" json:"payd_username_key"`
	PaydPasswordKey string `binding:"required" json:"payd_password_key"`
}

type RegisterUserResponse struct {
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (r *RabbitConn) HandleRegisterUser(req RegisterUserRequest) []byte {
	ctx, cancel := context.WithTimeout(context.Background(), 42*time.Second)
	defer cancel()

	user, err := r.UserRepository.CreateUser(ctx, repository.User{
		FullName:        req.FullName,
		Email:           req.Email,
		Password:        req.Password,
		PaydUsername:    req.PaydUsername,
		PaydAccountID:   req.PaydAccountID,
		PaydUsernameKey: req.PaydUsernameKey,
		PaydPasswordKey: req.PaydPasswordKey,
	})
	if err != nil {
		return errorRabbitMQResponse(
			pkg.Errorf(pkg.ErrorCode(err), "failed to create user: %v", pkg.ErrorMessage(err)),
		)
	}

	rsp := RegisterUserResponse{
		FullName:  user.FullName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	rspByte, rspErr := json.Marshal(rsp)
	if rspErr != nil {
		return errorRabbitMQResponse(
			pkg.Errorf(pkg.INTERNAL_ERROR, "failed to marshal response: %v", rspErr),
		)
	}

	return rspByte
}

type LoginUserRequest struct {
	Email    string `binding:"required" json:"email"`
	Password string `binding:"required" json:"password"`
}

type LoginUserResponse struct {
	AccessToken  string    `json:"access_token"`
	ExpirationAt time.Time `json:"expiration_at"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"created_at"`
}

func (r *RabbitConn) HandleLoginUser(req LoginUserRequest) []byte {
	ctx, cancel := context.WithTimeout(context.Background(), 42*time.Second)
	defer cancel()

	user, err := r.UserRepository.GetUser(ctx, req.Email)
	if err != nil {
		return errorRabbitMQResponse(
			pkg.Errorf(pkg.ErrorCode(err), "failed to login user: %v", pkg.ErrorMessage(err)),
		)
	}

	err = pkg.ComparePasswordAndHash(user.Password, req.Password)
	if err != nil {
		return errorRabbitMQResponse(
			pkg.Errorf(pkg.AUTHENTICATION_ERROR, "Error comparing passwords: %v", err),
		)
	}

	accessToken, err := r.Maker.CreateToken(user.Email, user.ID, r.Config.TOKEN_DURATION)
	if err != nil {
		return errorRabbitMQResponse(
			pkg.Errorf(pkg.AUTHENTICATION_ERROR, "Error creating token: %v", err),
		)
	}

	rsp := LoginUserResponse{
		AccessToken:  accessToken,
		ExpirationAt: time.Now().Add(r.Config.TOKEN_DURATION),
		FullName:     user.FullName,
		Email:        user.Email,
		CreatedAt:    user.CreatedAt,
	}

	rspByte, err := json.Marshal(rsp)
	if err != nil {
		return errorRabbitMQResponse(
			pkg.Errorf(pkg.INTERNAL_ERROR, "failed to marshal response: %v", err),
		)
	}

	return rspByte
}
