package main

import (
	"log"

	"github.com/EmilioCliff/payment-polling-app/payment-service/internal/http"
	"github.com/EmilioCliff/payment-polling-app/payment-service/internal/postgres"
	"github.com/EmilioCliff/payment-polling-app/payment-service/internal/rabbitmq"
	"github.com/EmilioCliff/payment-polling-app/payment-service/internal/workers"
	"github.com/EmilioCliff/payment-polling-app/payment-service/pkg"
	"github.com/EmilioCliff/payment-polling-service/shared-grpc/pb"
	"github.com/hibiken/asynq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	config, err := pkg.LoadConfig(".")
	if err != nil {
		log.Printf("error loading config: %s", err)

		return
	}

	store := postgres.NewStore(config)

	err = store.Start()
	if err != nil {
		log.Printf("error starting store: %s", err)

		return
	}

	transactionRepo := postgres.NewTransactionService(store)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.REDDIS_ADDR,
		DB:   1,
	}

	distributor := workers.NewRedisTaskDistributor(&redisOpt)

	processor, err := workers.NewRedisTaskProcessor(&redisOpt, config)
	if err != nil {
		log.Printf("error starting processor: %s", err)

		return
	}

	clientConn, err := grpc.NewClient(config.AUTH_GRPC_URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error connecting to auth grpc server: %s", err)

		return
	}
	defer clientConn.Close()

	client := pb.NewAuthenticationServiceClient(clientConn)

	rabbit := rabbitmq.NewRabbitConn(config, client)

	err = rabbit.ConnectToRabbit()
	if err != nil {
		log.Printf("error connecting to rabbit: %s", err)

		return
	}

	server, err := http.NewHttpServer()
	if err != nil {
		log.Printf("error starting server: %s", err)

		return
	}

	processor.TransactionRepository = transactionRepo

	rabbit.TransactionRepository = transactionRepo
	rabbit.Distributor = distributor

	server.TransactionRepository = transactionRepo

	go func() {
		processor.Start()
	}()

	go func() {
		rabbit.SetConsumer([]string{"payments.initiate_payment", "payments.poll_payments"})
	}()

	log.Println("Starting server on port", config.HTTP_PORT)

	err = server.Start(config.HTTP_PORT)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
