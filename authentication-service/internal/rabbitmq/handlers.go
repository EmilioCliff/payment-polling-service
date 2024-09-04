package rabbitmq

import (
	"context"
	"encoding/json"
	"time"

	"github.com/EmilioCliff/payment-polling-app/authentication-service/internal/postgres"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/pkg"
)

func (r *RabbitConn) handleRegisterUser(req postgres.RegisterUserRequest) []byte {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	rsp, err := r.store.RegisterUser(ctx, req)
	if err != nil {
		return r.errorRabbitMQResponse(pkg.Errorf(pkg.INTERNAL_ERROR, "failed to create user: %v", err))
	}

	rspByte, rspErr := json.Marshal(rsp)
	if rspErr != nil {
		return r.errorRabbitMQResponse(pkg.Errorf(pkg.INTERNAL_ERROR, "failed to marshal response: %v", rspErr))
	}

	return rspByte
}

func (r *RabbitConn) handleLoginUser(req postgres.LoginUserRequest) []byte {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	rsp, err := r.store.LoginUser(ctx, req)
	if err != nil {
		return r.errorRabbitMQResponse(pkg.Errorf(pkg.INTERNAL_ERROR, "failed to login user: %v", err))
	}

	rspByte, rspErr := json.Marshal(rsp)
	if rspErr != nil {
		return r.errorRabbitMQResponse(pkg.Errorf(pkg.INTERNAL_ERROR, "failed to marshal response: %v", rspErr))
	}

	return rspByte
}
