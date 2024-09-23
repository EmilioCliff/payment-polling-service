package main

import (
	"log"

	"github.com/EmilioCliff/payment-polling-app/authentication-service/internal/handlers/Grpc"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/internal/handlers/http"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/internal/handlers/rabbitmq"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/internal/postgres"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/pkg"
)

func main() {
	config, err := pkg.LoadConfig(".")
	if err != nil {
		log.Fatalf("failed to load config file: %v", err)
	}

	maker, err := pkg.NewJWTMaker(config.PRIVATE_KEY_PATH, config.PUBLIC_KEY_PATH)
	if err != nil {
		log.Fatalf("Failed to create token maker: %v", err)
	}

	db := postgres.NewStore(config)

	// connects to the db and run migrations
	err = db.Start()
	if err != nil {
		log.Fatalf("Failed to start postgres db: %v", err)
	}

	// create an instance of the user repository
	userRepository := postgres.NewUserService(db)

	grpcServer := Grpc.NewGRPCServer(config, *maker)
	grpcServer.UserRepository = userRepository

	rabbitConn := rabbitmq.NewRabbitConn(config, *maker)
	rabbitConn.UserRepository = userRepository

	// connects to rabbitmq
	if err = rabbitConn.ConnectToRabbit(); err != nil {
		log.Fatalf("Failed to create new rabbit conn: %v", err)
	}

	httpServer := http.NewHTTPServer(config, *maker)
	httpServer.UserRepository = userRepository

	go func() {
		log.Printf("Starting Authentication Grpc Server on port: %s", config.GRPC_PORT)

		if err = grpcServer.Start(config.GRPC_PORT); err != nil {
			log.Fatalf("Failed to create new gRPC server instance: %v", err)
		}
	}()

	go func() {
		log.Println("Starting Authentication Rabbit Consumer...")

		// declares the auth queue, binds topics, and starts consuming messages
		if err = rabbitConn.SetConsumer([]string{
			"authentication.register_user",
			"authentication.login_user",
		}); err != nil {
			log.Fatalf("Failed to create new rabbit conn: %v", err)
		}
	}()

	log.Printf("Starting Authentication Http Server on port: %s", config.HTTP_PORT)

	err = httpServer.Start()
	if err != nil {
		log.Fatalf("Failed to start http server instance: %v", err)
	}
}
