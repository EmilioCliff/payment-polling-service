package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/EmilioCliff/payment-polling-app/gateway-service/api"
	"github.com/EmilioCliff/payment-polling-app/gateway-service/utils"
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

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Failed to create gateway RabbitMQ channel: %s", err)
		return
	}
	defer ch.Close()

	server, err := api.NewServer(ch, config)
	if err != nil {
		log.Printf("Failed to start new server instance to db: %s", err)
		return
	}

	go server.SetConsumer([]string{"gateway.initiate_payment", "gateway.poll_payments", "gateway.register_user", "gateway.login_user"})

	log.Printf("Starting Gateway Server at port: %s", os.Getenv("PORT"))
	server.Start(fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT")))
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
