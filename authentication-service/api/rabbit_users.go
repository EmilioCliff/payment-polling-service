package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

var (
	ErrorUnmarshalingData = errors.New("error unmarshaling data")
	ErrorMarshalingData   = errors.New("error marshaling data")
)

func (server *Server) rabbitRegisterUser(payload registerUserRequest) []byte {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	rspData, errRsp := server.registerUserGeneral(payload, ctx)
	if !errRsp.Status {
		return server.errorRabbitMQResponse(errRsp.Error, errRsp.Message, errRsp.HttpStatus)
	}

	rsp, err := json.Marshal(rspData)
	if err != nil {
		return server.errorRabbitMQResponse(ErrorMarshalingData, "error marshaling response", http.StatusInternalServerError)
	}

	return rsp
}

func (server *Server) rabbitLoginUser(data loginUserRequest) []byte {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	rspData, errRsp := server.loginUserGeneral(data, ctx)
	if !errRsp.Status {
		return server.errorRabbitMQResponse(errRsp.Error, errRsp.Message, errRsp.HttpStatus)
	}

	rsp, err := json.Marshal(rspData)
	if err != nil {
		return server.errorRabbitMQResponse(ErrorMarshalingData, "error marshaling response", http.StatusInternalServerError)
	}

	return rsp
}
