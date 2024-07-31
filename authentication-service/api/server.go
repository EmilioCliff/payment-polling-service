package api

import (
	db "github.com/EmilioCliff/payment-polling-app/authentication-service/db/sqlc"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/pb"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/utils"
	"github.com/gin-gonic/gin"
)

type Server struct {
	pb.UnimplementedAuthenticationServiceServer
	router *gin.Engine
	config utils.Config
	store  *db.Queries
	maker  utils.JWTMaker
}

func NewServer(config utils.Config, store *db.Queries, maker utils.JWTMaker) (*Server, error) {
	server := &Server{
		config: config,
		store:  store,
		maker:  maker,
	}
	server.setRoutes()

	return server, nil
}

func (server *Server) setRoutes() {
	router := gin.Default()

	// auth := router.Group("/").Use(authenticationMiddleware(server.maker))

	router.POST("/auth/register", server.registerUser)
	router.POST("/auth/login", server.loginUser)

	server.router = router
}

func (server *Server) Start(addr string) error {
	return server.router.Run(addr)
}
