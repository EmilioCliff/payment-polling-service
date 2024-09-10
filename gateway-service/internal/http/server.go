package http

import (
	"fmt"

	_ "github.com/EmilioCliff/payment-polling-app/gateway-service/docs"
	"github.com/EmilioCliff/payment-polling-app/gateway-service/internal/gRPC"
	"github.com/EmilioCliff/payment-polling-app/gateway-service/internal/rabbitmq"
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

	GrpcClient   *gRPC.GrpcClient
	RabbitClient *rabbitmq.RabbitHandler
}

var _ services.HttpInterface = (*HttpServer)(nil)

func NewHttpServer(config pkg.Config, ch *amqp.Channel) (*HttpServer, error) {
	rpcClient, err := gRPC.NewGrpcClient()
	if err != nil {
		return nil, fmt.Errorf("Failed to start new gRPC client: %v", err)
	}

	rabbitHandler, err := rabbitmq.NewRabbitHandler(ch, config)
	if err != nil {
		return nil, fmt.Errorf("Failed to start rabbit handler: %v", err)
	}

	server := &HttpServer{
		config:       config,
		GrpcClient:   rpcClient,
		RabbitClient: rabbitHandler,
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
