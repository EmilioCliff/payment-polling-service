package main

import (
	"log"

	"github.com/EmilioCliff/payment-polling-app/authentication-service/internal/handlers/gRPC"
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

	db := postgres.NewStore(config, *maker)
	err = db.Start()
	if err != nil {
		log.Fatalf("Failed to start postgres db: %v", err)
	}

	userRepository := postgres.NewUserService(db)

	grpcServer := gRPC.NewGRPCServer(config, *maker)

	grpcServer.UserRepository = userRepository
	go func() {
		log.Printf("Starting Authentication Grpc Server on port: %s", config.GRPC_PORT)
		if err = grpcServer.Start(config.GRPC_PORT); err != nil {
			log.Fatalf("Failed to create new gRPC server instance: %v", err)
		}
	}()

	rabbitConn := rabbitmq.NewRabbitConn(config, *maker)

	if err = rabbitConn.ConnectToRabbit(); err != nil {
		log.Fatalf("Failed to create new rabbit conn: %v", err)
	}

	rabbitConn.UserRepository = userRepository
	go func() {
		log.Println("Starting Authentication Rabbit Consumer...")
		if err = rabbitConn.SetConsumer([]string{"authentication.register_user", "authentication.login_user"}); err != nil {
			log.Fatalf("Failed to create new rabbit conn: %v", err)
		}
	}()

	httpServer := http.NewHTTPServer(config, *maker)
	if err != nil {
		log.Fatalf("Failed to create new http server instance: %v", err)
	}

	httpServer.UserRepository = userRepository

	log.Printf("Starting Authentication Http Server on port: %s", config.HTTP_PORT)
	err = httpServer.Start()
	if err != nil {
		log.Fatalf("Failed to start http server instance: %v", err)
	}
}
