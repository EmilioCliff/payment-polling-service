package api

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
)

func (server *Server) errorResponse(err error, msg string) gin.H {
	return gin.H{
		"message": fmt.Errorf(msg+": %w", err).Error(),
	}
}

type errorRabbitResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (server *Server) errorRabbitMQResponse(err error, msg string, status int) []byte {
	errRsp := errorRabbitResponse{
		Status:  status,
		Message: fmt.Errorf(msg+": %w", err).Error(),
	}

	jsonResponse, err := json.Marshal(errRsp)
	if err != nil {
		log.Printf("Failed to marshal error response: %v", err)
		return []byte(`{"status": false}`)
	}

	return jsonResponse
}

type generalErrorResponse struct {
	Status     bool
	Message    string
	Error      error
	HttpStatus int
	GrpcCode   codes.Code
}

func errorHelper(msg string, err error, status int, code codes.Code) generalErrorResponse {
	return generalErrorResponse{
		Status:     false,
		Message:    msg,
		Error:      err,
		HttpStatus: status,
		GrpcCode:   code,
	}
}
