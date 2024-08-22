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
	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
)

func maini() {
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
	defer postgresConn.Close()

	// conn, err := amqp.Dial(config.RABBITMQ_URL)
	conn, err := connectToRabit(config.RABBITMQ_URL)
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ: %s", err)
		return
	}
	defer conn.Close()

	store := db.New(postgresConn)

	app := app.NewApp(conn, config, store)

	app.SetConsumer([]string{"payments.initiate_payment", "payments.poll_payments"})
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
