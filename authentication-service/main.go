package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"net"
	"time"

	"github.com/EmilioCliff/payment-polling-app/authentication-service/api"
	db "github.com/EmilioCliff/payment-polling-app/authentication-service/db/sqlc"

	// "github.com/EmilioCliff/payment-polling-app/authentication-service/pb"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/utils"

	pb "github.com/EmilioCliff/payment-polling-service/shared-grpc/pb"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Printf("Failed to load config files: %s", err)
		return
	}

	postgresConn, err := pgxpool.New(context.Background(), config.DB_URL)
	if err != nil {
		log.Printf("Failed to connect to db: %s", err)
		return
	}

	rabbitConn, err := connectToRabit(config.RABBITMQ_URL)
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ: %s", err)
		return
	}
	defer rabbitConn.Close()

	migration, err := migrate.New("file://db/migrations", config.DB_URL)
	if err != nil {
		log.Printf("Failed to load migration: %s", err)
		return
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Printf("Failed to run migrate up: %s", err)
		return
	}

	store := db.New(postgresConn)

	maker, err := utils.NewJWTMaker(config.PRIVATE_KEY_PATH, config.PUBLIC_KEY_PATH)
	if err != nil {
		log.Printf("Failed to create token maker: %s", err)
		return
	}

	server, err := api.NewServer(config, store, *maker)
	if err != nil {
		log.Printf("Failed to start new server instance to db: %s", err)
		return
	}

	go runGinServer(config.HTTP_PORT, server)

	go runRabbitMQConsumer(rabbitConn, server)

	rungRPCServer(config.GRPC_PORT, server)
}

func runGinServer(httpPort string, server *api.Server) {
	log.Printf("Starting Authentication Server at port: %s", httpPort)
	server.Start(httpPort)
}

func runRabbitMQConsumer(rabbitConn *amqp.Connection, server *api.Server) {
	err := server.SetConsumer([]string{"authentication.register_user", "authentication.login_user"}, rabbitConn)
	if err != nil {
		log.Printf("Failed to start rabbitMQ consumer: %s", err)
		return
	}
}

func rungRPCServer(grpcPort string, server *api.Server) {
	// rb.RegisterAuthenticationServiceServer(grpcPort, server)
	grpcServer := grpc.NewServer()

	pb.RegisterAuthenticationServiceServer(grpcServer, server)

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Printf("Failed to start grpc server on port: %s", err)
		return
	}

	log.Printf("Starting gRPC server on port: %s", grpcPort)
	if err = grpcServer.Serve(listener); err != nil {
		log.Printf("Failed to start gRPC server: %s", err)
		return
	}
}

func connectToRabit(uri string) (*amqp.Connection, error) {
	count := 0
	rollOff := 1 * time.Second
	var err error
	var connection *amqp.Connection
	for {
		connection, err = amqp.Dial(uri)
		if err != nil {
			log.Println("failed to connect to rabbitmq", err)
			if count > 12 {
				return nil, err
			}
			count++
			rollOff = time.Duration(math.Pow(float64(count), 2)) * time.Second
			time.Sleep(rollOff)
			continue
		}
		fmt.Println("Connected to rabbitmq")
		break
	}

	return connection, nil
}
