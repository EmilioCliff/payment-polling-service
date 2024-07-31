package api

import (
	"github.com/EmilioCliff/payment-polling-app/gateway-service/pb"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	router        *gin.Engine
	authgRPClient pb.AuthenticationServiceClient
}

func NewServer() (*Server, error) {
	conn, err := grpc.NewClient("authApp:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewAuthenticationServiceClient(conn)

	server := &Server{
		authgRPClient: client,
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
