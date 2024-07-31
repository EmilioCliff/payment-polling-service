package api

import (
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
}

func NewServer() *Server {
	server := &Server{}
	server.setRoutes()

	return server
}

func (server *Server) setRoutes() {
	router := gin.Default()

	// router.Use(CORSMiddleware())

	router.POST("/register", server.registerUser)
	router.POST("/login", server.loginUser)
	router.POST("/payments/initiate", server.initiatePayment)
	router.GET("/payments/status/:id", server.paymentStatus)

	server.router = router
}

func (server *Server) Start(addr string) error {
	return server.router.Run(addr)
}
