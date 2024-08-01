package app

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

func (app *App) SetConsumer(topics []string) error {
	ch, err := app.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		app.config.EXCH, // name
		"topic",         // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(
		app.config.QUEUE_NAME, // name
		false,                 // durable
		false,                 // delete when unused
		false,                 // exclusive
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		return err
	}

	for _, topic := range topics {
		if err := ch.QueueBind(
			q.Name,          // queue name
			topic,           // routing key
			app.config.EXCH, // exchange
			false,
			nil,
		); err != nil {
			return err
		}
	}

	messages, err := ch.Consume(
		q.Name,            // queue
		"payment_service", // consumer
		false,             // auto ack
		false,             // exclusive
		false,             // no local
		false,             // no wait
		nil,               // args
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
				log.Printf("failed to unmarshal message: %s", err)
				msg.Nack(false, true)
				return
			}

			response := distributeTask(payload)

			log.Printf("Message acknowledged: %v", msg.DeliveryTag)
			msg.Ack(false)

			err = ch.PublishWithContext(ctx,
				app.config.EXCH, // exchange
				msg.ReplyTo,     // routing key
				false,           // mandatory
				false,           // immediate
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
						app.config.EXCH, // exchange
						msg.ReplyTo,     // routing key
						false,           // mandatory
						false,           // immediate
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
			log.Println("response sent successfully")
		}
	}()

	log.Println("listening to messages in payment service")

	<-forever

	return nil
}

func distributeTask(payload Payload) []byte {
	switch payload.Name {
	case "initiate_payment":
		return initiatePayment(payload.Data.(map[string]interface{}))
	default:
		return nil
	}
}
