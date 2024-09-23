package http

import (
	"net/http"

	"github.com/EmilioCliff/payment-polling-app/authentication-service/internal/repository"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/pkg"
	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	router *gin.Engine
	config pkg.Config
	maker  pkg.JWTMaker

	UserRepository repository.UserRepository
}

func NewHTTPServer(config pkg.Config, tokenMaker pkg.JWTMaker) *HTTPServer {
	s := &HTTPServer{
		config: config,
		maker:  tokenMaker,
	}

	s.setRoutes()

	return s
}

func (s *HTTPServer) setRoutes() {
	r := gin.Default()

	r.GET("/healthcheck", s.handleHealthCheck)
	r.POST("/auth/register", s.handleRegisterUser)
	r.POST("/auth/login", s.handleLoginUser)

	s.router = r
}

func (s *HTTPServer) handleHealthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

func (s *HTTPServer) Start() error {
	return s.router.Run(s.config.HTTP_PORT)
}
