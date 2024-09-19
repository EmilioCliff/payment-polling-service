package http

import (
	"fmt"
	"net/http"

	"github.com/EmilioCliff/payment-polling-app/payment-service/internal/postgres"
	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	router *gin.Engine
	store  postgres.Store
}

func NewHttpServer() (*HttpServer, error) {
	server := &HttpServer{}

	store, err := postgres.GetStore()
	if err != nil {
		return nil, fmt.Errorf("failed to create new store in server: %v", err)
	}

	server.store = store

	server.setRoutes()

	return server, nil
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
