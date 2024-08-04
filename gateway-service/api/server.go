package api

import (
	"sync"

	"github.com/EmilioCliff/payment-polling-app/gateway-service/pb"
	"github.com/EmilioCliff/payment-polling-app/gateway-service/utils"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	router        *gin.Engine
	authgRPClient pb.AuthenticationServiceClient
	amqpChannel   *amqp.Channel
	config        utils.Config
	responseMap   sync.Map
}

func NewServer(amqpChannel *amqp.Channel, config utils.Config) (*Server, error) {
	gRPCconn, err := grpc.NewClient("authApp:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewAuthenticationServiceClient(gRPCconn)

	server := &Server{
		authgRPClient: client,
		amqpChannel:   amqpChannel,
		config:        config,
	}

	server.setRoutes()

	return server, nil
}

func (server *Server) setRoutes() {
	router := gin.Default()

	router.Use(CORSMiddleware())

	router.POST("/register", server.registerUser)
	router.POST("/login", server.loginUser)
	router.POST("/payments/initiate", server.initiatePayment)
	router.GET("/payments/status/:id", server.paymentStatus)

	server.router = router
}

func (server *Server) Start(addr string) error {
	return server.router.Run(addr)
}
