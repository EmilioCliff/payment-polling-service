package http

import (
	"net/http"

	"github.com/EmilioCliff/payment-polling-app/gateway-service/internal/services"
	"github.com/EmilioCliff/payment-polling-app/gateway-service/pkg"
	"github.com/gin-gonic/gin"
)

func (s *HttpServer) handleRegisterUser(ctx *gin.Context) {
	var req services.RegisterUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse("Invalid request", http.StatusBadRequest))

		return
	}

	// we can change the communication channel here ie
	// statusCode, rsp := s.RabbitService.RegisterUserViaRabbit(req)
	// statusCode, rsp := s.HTTPService.RegisterUserViaHttp(req)
	statusCode, rsp := s.GRPCService.RegisterUserViagRPC(req)
	if statusCode != http.StatusOK {
		ctx.JSON(statusCode, pkg.ErrorResponse(rsp.Message, rsp.StatusCode))

		return
	}

	ctx.JSON(statusCode, rsp)
}

func (s *HttpServer) handleLoginUser(ctx *gin.Context) {
	var req services.LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse("Invalid request", http.StatusBadRequest))

		return
	}

	// we can change the communication channel here ie
	// statusCode, rsp := s.GrpcClient.LoginUserViagRPC(req)
	// statuSCode, rsp := s.RabbitService.LoginUserViaRabbit(req)
	statusCode, rsp := s.HTTPService.LoginUserViaHttp(req)
	if statusCode != http.StatusOK {
		ctx.JSON(statusCode, pkg.ErrorResponse(rsp.Message, rsp.StatusCode))

		return
	}

	ctx.JSON(statusCode, rsp)
}

func (s *HttpServer) handleInitiatePayment(ctx *gin.Context) {
	var req services.InitiatePaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse("Invalid request", http.StatusBadRequest))

		return
	}

	// implemented only RabbitMQ communication channel with payment service.
	statusCode, rsp := s.RabbitService.InitiatePaymentViaRabbit(req)
	if statusCode != http.StatusOK {
		ctx.JSON(statusCode, pkg.ErrorResponse(rsp.Message, rsp.StatusCode))

		return
	}

	ctx.JSON(statusCode, rsp)
}

func (s *HttpServer) handlePaymentPolling(ctx *gin.Context) {
	var req services.PollingTransactionRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse("Invalid request", http.StatusBadRequest))

		return
	}

	value, exists := ctx.Get(authorizationPayloadKey)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, pkg.ErrorResponse("Missing token payload", http.StatusUnauthorized))

		return
	}

	payload, ok := value.(*pkg.Payload)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "type assertion failed"})

		return
	}

	// implemented only RabbitMQ communication channel with payment service.
	statusCode, rsp := s.RabbitService.PollTransactionViaRabbit(req, payload.UserID)
	if statusCode != http.StatusOK {
		ctx.JSON(statusCode, pkg.ErrorResponse(rsp.Message, rsp.StatusCode))

		return
	}

	ctx.JSON(statusCode, rsp)
}
