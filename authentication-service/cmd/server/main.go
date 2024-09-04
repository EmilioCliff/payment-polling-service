package main

import (
	"log"
	"os"

	"github.com/EmilioCliff/payment-polling-app/authentication-service/internal/gRPC"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/internal/http"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/internal/rabbitmq"
)

func main() {
	rabbitConn, err := rabbitmq.NewRabbitConn()
	if err != nil {
		log.Printf("Failed to create new rabbit conn: %v", err)
		os.Exit(1)
	}

	httpServer, err := http.NewHttpServer()
	if err != nil {
		log.Printf("Failed to create new http server instance: %v", err)
		os.Exit(1)
	}

	grpcServer, err := gRPC.NewGrpcServer()
	if err != nil {
		log.Printf("Failed to create new gRPC server instance: %v", err)
		os.Exit(1)
	}

	log.Printf("Starting Authentication Grpc Server on port: %s", grpcServer.Addr)
	go grpcServer.Start()

	go func() {
		log.Printf("Starting Authentication Http Server on port: %s", httpServer.Addr)
		err = httpServer.Start()
		if err != nil {
			log.Printf("Failed to start new server: %v", err)
			os.Exit(1)
		}
	}()

	err = rabbitConn.SetConsumer([]string{"authentication.register_user", "authentication.login_user"})
	if err != nil {
		log.Printf("Failed to create new rabbit conn: %v", err)
		os.Exit(1)
	}
}
