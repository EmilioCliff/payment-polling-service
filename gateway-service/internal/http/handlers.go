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
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse(err, "Invalid request"))
		return
	}

	statusCode, rsp := s.RegisterUserViaHttp(req)
	ctx.JSON(statusCode, rsp)
}

func (s *HttpServer) handleLoginUser(ctx *gin.Context) {
	var req services.LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse(err, "Invalid request"))
		return
	}

	statusCode, rsp := s.GrpcClient.LoginUserViagRPC(req)
	ctx.JSON(statusCode, rsp)
}

func (s *HttpServer) handleInitiatePayment(ctx *gin.Context) {
	var req services.InitiatePaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse(err, "Invalid request"))
		return
	}

	statusCode, rsp := s.RabbitClient.InitiatePaymentViaRabbit(req)
	ctx.JSON(statusCode, rsp)
}

func (s *HttpServer) handlePaymentPolling(ctx *gin.Context) {
	var req services.PollingTransactionRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse(err, "Invalid request"))
		return
	}

	statusCode, rsp := s.RabbitClient.PollTransactionViaRabbit(req)
	ctx.JSON(statusCode, rsp)
}
