package main

import (
	"fmt"
	"log"

	"github.com/EmilioCliff/payment-polling-app/payment-service/internal/http"
	"github.com/EmilioCliff/payment-polling-app/payment-service/internal/rabbitmq"
	"github.com/EmilioCliff/payment-polling-app/payment-service/internal/workers"
	"github.com/EmilioCliff/payment-polling-app/payment-service/pkg"
	"github.com/hibiken/asynq"
)

func main() {
	config, err := pkg.LoadConfig(".")
	if err != nil {
		log.Printf("error loading config: %s", err)
		return
	}

	rabbit, err := rabbitmq.NewRabbitConn()
	if err != nil {
		log.Printf("error connecting to rabbit: %s", err)
		return
	}

	server, err := http.NewHttpServer()
	if err != nil {
		log.Printf("error starting server: %s", err)
		return
	}

	redisOpt := asynq.RedisClientOpt{
		Addr: config.REDDIS_ADDR,
		DB:   1,
	}

	processor, err := workers.NewRedisTaskProcessor(&redisOpt)
	if err != nil {
		log.Printf("error starting processor: %s", err)
		return
	}

	processorErrChan := make(chan error, 1)
	rabbitErrChan := make(chan error, 1)

	go func() {
		processorErrChan <- processor.Start()
	}()

	go func() {
		rabbitErrChan <- rabbit.SetConsumer([]string{"payments.initiate_payment", "payments.poll_payments"})
	}()

	select {
	case err := <-processorErrChan:
		log.Fatalf("Error starting processor: %v", err)
		close(processorErrChan)
		return
	case err := <-rabbitErrChan:
		log.Fatalf("Error starting rabbit consumer: %v", err)
		close(rabbitErrChan)
		return
	default:
		fmt.Println("Both processor and rabbit consumer started successfully.")
	}

	log.Println("Starting server: ", config.HTTP_PORT)
	err = server.Start(config.HTTP_PORT)
	if err != nil {
		log.Printf("error starting server: %s", err)
		return
	}
}
