package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/EmilioCliff/payment-polling-app/payment-service/app"
	db "github.com/EmilioCliff/payment-polling-app/payment-service/db/sqlc"
	"github.com/EmilioCliff/payment-polling-app/payment-service/utils"
	"github.com/EmilioCliff/payment-polling-app/payment-service/workers"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Printf("Failed to load config files: %s", err)
		return
	}

	// conn, err := amqp.Dial(config.RABBITMQ_URL)
	conn, err := connectToRabit(config.RABBITMQ_URL)
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ: %s", err)
		return
	}
	defer conn.Close()

	postgresConn, err := pgxpool.New(context.Background(), config.DB_URL)
	if err != nil {
		log.Printf("Failed to connect to db: %s", err)
		return
	}
	defer postgresConn.Close()

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

	app, err := app.NewApp(conn, config, store)
	if err != nil {
		log.Printf("Failed to create new app: %s", err)
		return
	}

	go startAsynqProcessor(store)
	go app.SetConsumer([]string{"payments.initiate_payment", "payments.poll_payments"})

	log.Println("starting payment service server on port: 3030")
	err = app.Start("0.0.0.0:3030")
	if err != nil {
		log.Printf("Failed to start app: %s", err)
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

func startAsynqProcessor(store *db.Queries) {
	redisOpt := asynq.RedisClientOpt{
		Addr: "redis:6379",
		DB:   1,
	}
	processor := workers.NewRedisTaskProcessor(store, &redisOpt)
	err := processor.Start()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("asynq processor started successfully")
}
