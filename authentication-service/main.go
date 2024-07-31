package main

import (
	"context"
	"log"
	"net"

	"github.com/EmilioCliff/payment-polling-app/authentication-service/api"
	db "github.com/EmilioCliff/payment-polling-app/authentication-service/db/sqlc"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/pb"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/utils"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// main is the entry point of the application.
//
// It initializes the configuration, establishes a connection to the database,
// creates a JWT token maker, creates a database store, and starts the
// authentication server.
//
// Parameters:
//
//	None.
//
// Returns:
//
//	None.
func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Printf("Failed to load config files: %s", err)
		return
	}

	conn, err := pgxpool.New(context.Background(), config.DB_URL)
	if err != nil {
		log.Printf("Failed to connect to db: %s", err)
		return
	}

	migration, err := migrate.New("file://db/migrations", config.DB_URL)
	if err != nil {
		log.Printf("Failed to load migration: %s", err)
		return
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Printf("Failed to run migrate up: %s", err)
		return
	}

	store := db.New(conn)

	maker, err := utils.NewJWTMaker(config.PRIVATE_KEY_PATH, config.PUBLIC_KEY_PATH)
	if err != nil {
		log.Printf("Failed to create token maker: %s", err)
		return
	}

	go runGinServer(config, store, maker)

	rungRPCServer(config, store, maker)
}

func runGinServer(config utils.Config, store *db.Queries, maker *utils.JWTMaker) {
	server, err := api.NewServer(config, store, *maker)
	if err != nil {
		log.Printf("Failed to start new server instance to db: %s", err)
		return
	}

	log.Printf("Starting Authentication Server at port: %s", config.HTTP_PORT)
	server.Start(config.HTTP_PORT)
}

func rungRPCServer(config utils.Config, store *db.Queries, maker *utils.JWTMaker) {
	grpcServer := grpc.NewServer()

	server, err := api.NewServer(config, store, *maker)
	if err != nil {
		log.Printf("Failed to start new server instance to db: %s", err)
		return
	}

	pb.RegisterAuthenticationServiceServer(grpcServer, server)

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPC_PORT)
	if err != nil {
		log.Printf("Failed to start grpc server on port: %s", err)
		return
	}

	log.Printf("Starting gRPC server on port: %s", config.GRPC_PORT)
	if err = grpcServer.Serve(listener); err != nil {
		log.Printf("Failed to start gRPC server: %s", err)
		return
	}
}
