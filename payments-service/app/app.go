package app

import (
	db "github.com/EmilioCliff/payment-polling-app/payment-service/db/sqlc"
	"github.com/EmilioCliff/payment-polling-app/payment-service/utils"
	pb "github.com/EmilioCliff/payment-polling-service/shared-grpc/pb"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	conn          *amqp.Connection
	config        utils.Config
	store         *db.Queries
	authgRPClient pb.AuthenticationServiceClient
}

func NewApp(conn *amqp.Connection, config utils.Config, store *db.Queries) (*App, error) {
	gRPCconn, err := grpc.NewClient("authApp:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewAuthenticationServiceClient(gRPCconn)

	app := App{
		config:        config,
		conn:          conn,
		store:         store,
		authgRPClient: client,
	}

	return &app, nil
}
