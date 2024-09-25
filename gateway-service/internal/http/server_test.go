package http

import (
	"crypto/rand"
	"crypto/rsa"

	"github.com/EmilioCliff/payment-polling-app/gateway-service/internal/mock"
	"github.com/EmilioCliff/payment-polling-app/gateway-service/pkg"
	"github.com/gin-gonic/gin"
)

type TestHttpServer struct {
	server *HttpServer

	GrpcService   mock.MockGrpcService
	HTTPService   mock.MockHttpService
	RabbitService mock.MockRabbitMQService
}

func NewTestHttpServer() *TestHttpServer {
	gin.SetMode(gin.TestMode)

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	publicKey := &privateKey.PublicKey

	s := &TestHttpServer{
		server: NewHttpServer(pkg.JWTMaker{
			PublicKey:  publicKey,
			PrivateKey: privateKey,
		}),
	}

	s.server.HTTPService = &s.HTTPService
	s.server.GRPCService = &s.GrpcService
	s.server.RabbitService = &s.RabbitService

	return s
}
