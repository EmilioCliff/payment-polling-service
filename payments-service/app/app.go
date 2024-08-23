package app

import (
	db "github.com/EmilioCliff/payment-polling-app/payment-service/db/sqlc"
	"github.com/EmilioCliff/payment-polling-app/payment-service/utils"
	"github.com/EmilioCliff/payment-polling-app/payment-service/workers"
	pb "github.com/EmilioCliff/payment-polling-service/shared-grpc/pb"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	conn          *amqp.Connection
	config        utils.Config
	store         *db.Queries
	authgRPClient pb.AuthenticationServiceClient
	router        *gin.Engine
	distributor   workers.TaskDistributor
}

func NewApp(conn *amqp.Connection, config utils.Config, store *db.Queries) (*App, error) {
	gRPCconn, err := grpc.NewClient("authApp:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	redisOpt := asynq.RedisClientOpt{
		Addr: "redis:6379",
		DB:   1,
	}

	distributor := workers.NewRedisTaskDistributor(redisOpt)

	client := pb.NewAuthenticationServiceClient(gRPCconn)

	app := App{
		config:        config,
		conn:          conn,
		store:         store,
		authgRPClient: client,
		distributor:   distributor,
	}

	app.setRoutes()

	return &app, nil
}

func (app *App) setRoutes() {
	r := gin.Default()

	r.POST("/callback/:id", app.callBack)

	app.router = r
}

func (app *App) Start(addr string) error {
	return app.router.Run(addr)
}
