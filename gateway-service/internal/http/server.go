package http

import (
	_ "github.com/EmilioCliff/payment-polling-app/gateway-service/docs"
	"github.com/EmilioCliff/payment-polling-app/gateway-service/internal/services"
	"github.com/EmilioCliff/payment-polling-app/gateway-service/pkg"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type HttpServer struct {
	router *gin.Engine
	config pkg.Config

	HTTPService services.HttpInterface
	RabbitService services.RabbitInterface
	GRPCService   services.GrpcInterface
}

func NewHttpServer(config pkg.Config, ch *amqp.Channel) (*HttpServer, error) {
	server := &HttpServer{
		config:       config,
	}

	server.setRoutes()

	return server, nil
}

func (s *HttpServer) setRoutes() {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.POST("/register", s.handleRegisterUser)
	r.POST("/login", s.handleLoginUser)                   // need authentication
	r.POST("/payments/initiate", s.handleInitiatePayment) // need authentication
	r.GET("/payments/status/:id", s.handlePaymentPolling) // need authentication

	s.router = r
}

func (s *HttpServer) Start(addr string) error {
	return s.router.Run(addr)
}
