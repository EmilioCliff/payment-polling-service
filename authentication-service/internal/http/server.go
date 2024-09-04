package http

import (
	"net/http"

	"github.com/EmilioCliff/payment-polling-app/authentication-service/internal/postgres"
	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	router *gin.Engine
	store  *postgres.Store

	Addr string
}

func NewHttpServer() (*HttpServer, error) {
	store, err := postgres.NewStore()
	if err != nil {
		return nil, err
	}

	server := &HttpServer{
		store: store,
		Addr:  store.HTTP_ADDR,
	}

	server.setRoutes()

	return server, nil
}

func (s *HttpServer) setRoutes() {
	r := gin.Default()

	r.GET("/healthcheck", s.handleHealthCheck)
	r.POST("/auth/register", s.handleRegisterUser)
	r.POST("/auth/login", s.handleLoginUser)

	s.router = r
}

func (s *HttpServer) handleHealthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

func (s *HttpServer) Start() error {
	return s.router.Run(s.Addr)
}
