package rabbitmq_test

import (
	"crypto/rand"
	"crypto/rsa"
	"time"

	"github.com/EmilioCliff/payment-polling-app/authentication-service/internal/handlers/rabbitmq"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/internal/mock"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/pkg"
)

type TestRabbitConn struct {
	rabbitConn     *rabbitmq.RabbitConn
	UserRepository mock.MockUsersRepositry
}

func NewTestRabbitConn() *TestRabbitConn {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	publicKey := &privateKey.PublicKey

	r := TestRabbitConn{
		rabbitConn: rabbitmq.NewRabbitConn(
			pkg.Config{TOKEN_DURATION: time.Second},
			pkg.JWTMaker{PublicKey: publicKey, PrivateKey: privateKey},
		),
	}

	r.rabbitConn.UserRepository = &r.UserRepository

	return &r
}

// func TestRabbitConn_DistributeTask(t *testing.T) {
// 	r := NewTestRabbitConn()

// 	tests := []struct{
// 		name string
// 		payload rabbitmq.Payload

// 	}{}
// }
