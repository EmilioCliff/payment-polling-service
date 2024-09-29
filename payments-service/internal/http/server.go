package http

import (
	"net/http"

	"github.com/EmilioCliff/payment-polling-app/payment-service/internal/repository"
	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	router *gin.Engine

	TransactionRepository repository.TransactionRepository
}

func NewHttpServer() *HttpServer {
	server := &HttpServer{}

	server.setRoutes()

	return server
}

func (s *HttpServer) setRoutes() {
	r := gin.Default()

	r.GET("/healthcheck", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})
	r.POST("/transaction/:id", s.handleCallBack)

	s.router = r
}

func (s *HttpServer) Start(addr string) error {
	return s.router.Run(addr)
}
