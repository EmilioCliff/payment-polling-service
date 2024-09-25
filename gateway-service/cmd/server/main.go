package main

import (
	"log"

	_ "github.com/EmilioCliff/payment-polling-app/gateway-service/docs"
	"github.com/EmilioCliff/payment-polling-app/gateway-service/internal/gRPC"
	"github.com/EmilioCliff/payment-polling-app/gateway-service/internal/http"
	"github.com/EmilioCliff/payment-polling-app/gateway-service/internal/rabbitmq"
	"github.com/EmilioCliff/payment-polling-app/gateway-service/pkg"
)

// @title Payment Polling App
// @version 1.0
// @description Payment Polling App is an online polling payment gateway service.

// @host localhost:8080
func main() {
	config, err := pkg.LoadConfig(".")
	if err != nil {
		log.Printf("Failed to load config: %v", err)

		return
	}

	maker, err := pkg.NewJWTMaker(config.PRIVATE_KEY_PATH, config.PUBLIC_KEY_PATH)
	if err != nil {
		log.Fatalf("Failed to create token maker: %v", err)
	}

	conn, err := rabbitmq.ConnectToRabit(config.RABBITMQ_URL)
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ: %s", err)

		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Failed to create gateway RabbitMQ channel: %s", err)

		return
	}
	defer ch.Close()

	rabbitHandler := rabbitmq.NewRabbitService(ch, config)

	httpService := http.NewHTTPService(config)

	rpcClient := gRPC.NewGrpcService()

	err = rpcClient.Start(config.AUTH_GRPC_PORT)
	if err != nil {
		log.Fatalf("failed to get grpc client: %v", err)

		return
	}

	server := http.NewHttpServer(*maker)

	// injecting applications dependencies
	server.RabbitService = rabbitHandler
	server.HTTPService = httpService
	server.GRPCService = rpcClient

	readyCh := make(chan struct{}, 1)

	// declares the gateway queue, binds topics, and starts consuming messages
	go server.RabbitService.SetConsumer(
		[]string{
			"gateway.initiate_payment",
			"gateway.poll_payments",
			"gateway.register_user",
			"gateway.login_user",
		}, readyCh,
	)

	log.Printf("Starting server on port: %s", config.SERVER_ADDRESS)

	err = server.Start(config.SERVER_ADDRESS)
	if err != nil {
		log.Printf("Failed to start gateway service: %s", err)
	}

	log.Println("Shutting down gateway service...")
}

// func connectToRabit(uri string) (*amqp.Connection, error) {
// 	count := 0
// 	maxRetries := 12

// 	var err error

// 	var connection *amqp.Connection

// 	for {
// 		connection, err = amqp.Dial(uri)
// 		if err != nil {
// 			log.Println("failed to connect to rabbitmq", err)

// 			if count > maxRetries {
// 				return nil, err
// 			}

// 			count++
// 			rollOff := time.Duration(math.Pow(float64(count), 2)) * time.Second
// 			time.Sleep(rollOff)

// 			continue
// 		}

// 		log.Println("Connected to rabbitmq")

// 		break
// 	}

// 	return connection, nil
// }
