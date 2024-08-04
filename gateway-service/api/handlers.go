package api

import (
	"github.com/gin-gonic/gin"
)

func (server *Server) registerUser(ctx *gin.Context) {
	server.registerUserViaRabbit(ctx)
	// registerUserViagRPC(ctx, server)
	// registerUserViaHTTP(ctx, server)
}

func (server *Server) loginUser(ctx *gin.Context) {
	loginUserViaRabbit(ctx, server)
	// loginUserViagRPC(ctx, server)
	// loginUserViaHTTP(ctx, server)
}

func (server *Server) initiatePayment(ctx *gin.Context) {
	server.initiatePaymentViaRabbitMQ(ctx)
}

func (server *Server) paymentStatus(ctx *gin.Context) {

}
