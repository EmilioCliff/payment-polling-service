package api

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Payload struct {
	Name string `json:"name"`
	Data any    `json:"data"`
}

func (server *Server) SetConsumer(topics []string, rabbitConn *amqp.Connection) error {
	server.rabbitConn = rabbitConn

	ch, err := server.rabbitConn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		server.config.AUTH_QUEUE_NAME, // name
		false,                         // durable
		false,                         // delete when unused
		false,                         // exclusive
		false,                         // no-wait
		nil,                           // arguments
	)
	if err != nil {
		return err
	}

	for _, topic := range topics {
		if err := ch.QueueBind(
			q.Name,             // queue name
			topic,              // routing key
			server.config.EXCH, // exchange
			false,
			nil,
		); err != nil {
			return err
		}
	}

	messages, err := ch.Consume(
		q.Name,                   // queue
		"authentication_service", // consumer
		false,                    // auto ack
		false,                    // exclusive
		false,                    // no local
		false,                    // no wait
		nil,                      // args
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		for msg := range messages {
			var payload Payload
			err := json.Unmarshal(msg.Body, &payload)
			if err != nil {
				log.Printf("failed to unmarshal message in auth rabbit: %s", err)
				msg.Nack(false, true)
				return
			}

			response := server.distributeTask(payload)

			log.Printf("Message acknowledged from auth service: %v", msg.DeliveryTag)
			msg.Ack(false)

			err = ch.PublishWithContext(ctx,
				server.config.EXCH, // exchange
				msg.ReplyTo,        // routing key
				false,              // mandatory
				false,              // immediate
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: msg.CorrelationId,
					Body:          response,
				},
			)

			if err != nil {
				log.Printf("failed to send response: %s", err)
				count := 0
				for {
					err = ch.PublishWithContext(ctx,
						server.config.EXCH, // exchange
						msg.ReplyTo,        // routing key
						false,              // mandatory
						false,              // immediate
						amqp.Publishing{
							ContentType:   "text/plain",
							CorrelationId: msg.CorrelationId,
							Body:          response,
						},
					)
					if err == nil {
						break
					} else {
						count++
						if count > 5 {
							// log to failed resposnes
							log.Printf("failed to send response: %s", err)
							return
						}
					}
				}
			}
			log.Println("response sent successfully from auth service")
		}
	}()

	log.Println("listening to messages in authentication service")

	<-forever

	return nil
}

func (server *Server) distributeTask(payload Payload) []byte {
	switch payload.Name {
	case "register_user":
		log.Println("registering user in auth seerver ditributor")
		return server.rabbitRegisterUser(payload.Data.(map[string]interface{}))
	case "login_user":
		return server.rabbitLoginUser(payload.Data.(map[string]interface{}))
	default:
		return nil
	}
}
