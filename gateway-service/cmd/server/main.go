package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/EmilioCliff/payment-polling-app/gateway-service/internal/http"
	"github.com/EmilioCliff/payment-polling-app/gateway-service/pkg"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	config, err := pkg.LoadConfig(".")
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		return
	}

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

	server, err := http.NewHttpServer(config, ch)
	if err != nil {
		log.Printf("Failed to create new httpServer instance: %s", err)
		return
	}

	go server.RabbitClient.SetConsumer([]string{"gateway.initiate_payment", "gateway.poll_payments", "gateway.register_user", "gateway.login_user"})

	log.Printf("Starting server on port: %s", config.SERVER_ADDRESS)
	err = server.Start(config.SERVER_ADDRESS)
	if err != nil {
		log.Printf("Failed to start gateway service: %s", err)
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
