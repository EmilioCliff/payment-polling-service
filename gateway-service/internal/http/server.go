package http

import (
	"log"

	"github.com/EmilioCliff/payment-polling-app/gateway-service/internal/services"
	"github.com/EmilioCliff/payment-polling-app/gateway-service/pkg"
	"github.com/gin-gonic/gin"
	"github.com/rakyll/statik/fs"
)

type HttpServer struct {
	router *gin.Engine
	maker  pkg.JWTMaker

	HTTPService   services.HttpInterface
	RabbitService services.RabbitInterface
	GRPCService   services.GrpcInterface
}

func NewHttpServer(maker pkg.JWTMaker) *HttpServer {
	server := &HttpServer{
		maker: maker,
	}

	server.setRoutes()

	return server
}

func (s *HttpServer) setRoutes() {
	r := gin.Default()

	auth := r.Group("/").Use(authenticationMiddleware(s.maker)) // requires access token

	statikFs, err := fs.New()
	if err != nil {
		log.Fatal("cannot create statik fs")
	}

	r.StaticFS("/swagger", statikFs)

	r.POST("/register", s.handleRegisterUser)
	r.POST("/login", s.handleLoginUser)
	auth.POST("/payments/initiate", s.handleInitiatePayment)
	auth.GET("/payments/status/:id", s.handlePaymentPolling)

	s.router = r
}

func (s *HttpServer) Start(addr string) error {
	return s.router.Run(addr)
}
