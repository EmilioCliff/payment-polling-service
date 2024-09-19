package http

import (
	"crypto/rand"
	"crypto/rsa"
	"time"

	"github.com/EmilioCliff/payment-polling-app/authentication-service/internal/mock"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/pkg"
	"github.com/gin-gonic/gin"
)

type TestHTTPServer struct {
	server  *HTTPServer
	UserRepository mock.MockUsersRepositry
}

func NewTestHTTPServer() *TestHTTPServer {
	gin.SetMode(gin.TestMode)

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
    if err != nil {
        panic(err)
    }

    publicKey := &privateKey.PublicKey


	s := &TestHTTPServer{
		server: NewHTTPServer(pkg.Config{
			HTTP_PORT: "0.0.0.0:5050",
			TOKEN_DURATION: time.Second,
		}, pkg.JWTMaker{
			PublicKey: publicKey,
			PrivateKey: privateKey,
		}),
	}

	s.server.UserRepository = &s.UserRepository

	return s
}

// func TestHTTPServer_Start(t *testing.T) {
// 	s := NewTestHTTPServer()

// 	err := s.server.Start()
// 	require.NoError(t, err)
// }