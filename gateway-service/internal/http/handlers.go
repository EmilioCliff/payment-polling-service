package http

import (
	"net/http"

	_ "github.com/EmilioCliff/payment-polling-app/gateway-service/docs"
	"github.com/EmilioCliff/payment-polling-app/gateway-service/internal/services"
	"github.com/EmilioCliff/payment-polling-app/gateway-service/pkg"
	"github.com/gin-gonic/gin"
)

// RegisterUser godoc
// @Summary      Register a user
// @Description  Registers a new user in the database.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body services.RegisterUserRequest true "users details"
// @Success      200 {object} services.RegisterUserResponse{}
// @Router       /register [post]
func (s *HttpServer) handleRegisterUser(ctx *gin.Context) {
	var req services.RegisterUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse(err, "Invalid request"))
		return
	}

	statusCode, rsp := s.GRPCService.RegisterUserViagRPC(req)
	ctx.JSON(statusCode, rsp)
}

// LoginUser godoc
// @Summary      Login a user
// @Description  Logs in a user with credentials.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body services.LoginUserRequest true "users credetials"
// @Success      200 {object} services.LoginUserResponse{}
// @Router       /login [post]
func (s *HttpServer) handleLoginUser(ctx *gin.Context) {
	var req services.LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse(err, "Invalid request"))
		return
	}

	// statusCode, rsp := s.GrpcClient.LoginUserViagRPC(req)
	statusCode, rsp := s.HTTPService.LoginUserViaHttp(req)
	ctx.JSON(statusCode, rsp)
}

// InitiatePayment godoc
// @Summary      Initiate a payment
// @Description  Initiates a payment transaction.
// @Tags         payments
// @Accept       json
// @Produce      json
// @Param        request body services.InitiatePaymentRequest true "payment details"
// @Success      200 {object} services.InitiatePaymentResponse{}
// @Router       /initiate-payment [post]
func (s *HttpServer) handleInitiatePayment(ctx *gin.Context) {
	var req services.InitiatePaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse(err, "Invalid request"))
		return
	}

	statusCode, rsp := s.RabbitService.InitiatePaymentViaRabbit(req)
	ctx.JSON(statusCode, rsp)
}

// PollPayment godoc
// @Summary      Poll a payment
// @Description  Polls the status of a payment transaction.
// @Tags         payments
// @Accept       json
// @Produce      json
// @Param        id path string true "Transaction ID"
// @Success      200 {object} services.PollingTransactionResponse{}
// @Router       /poll-payment/{id} [get]
func (s *HttpServer) handlePaymentPolling(ctx *gin.Context) {
	var req services.PollingTransactionRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse(err, "Invalid request"))
		return
	}

	statusCode, rsp := s.RabbitService.PollTransactionViaRabbit(req)
	ctx.JSON(statusCode, rsp)
}
