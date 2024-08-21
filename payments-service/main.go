package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/EmilioCliff/payment-polling-app/payment-service/app"
	"github.com/EmilioCliff/payment-polling-app/payment-service/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

func maini() {
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

	app := app.NewApp(conn, config)

	app.SetConsumer([]string{"payments.initiate_payment", "poll_payments"})
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
