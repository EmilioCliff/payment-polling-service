package api

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	INITIATE_PAYMENT = "gateway.initiate_payment"
)

func (server *Server) SetConsumer(topics []string) error {
	q, err := server.amqpChannel.QueueDeclare(
		server.config.EXCLUSIVE_QUEUE_NAME, // name
		false,                              // durable
		false,                              // delete when unused
		false,                              // exclusive
		false,                              // no-wait
		nil,                                // arguments
	)
	if err != nil {
		return err
	}

	for _, topic := range topics {
		if err := server.amqpChannel.QueueBind(
			q.Name,             // queue name
			topic,              // routing key
			server.config.EXCH, // exchange
			false,
			nil,
		); err != nil {
			return err
		}
	}

	messages, err := server.amqpChannel.Consume(
		q.Name,            // queue
		"gateway_service", // consumer
		true,              // auto ack
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
		for msg := range messages {
			if ch, ok := server.responseMap.Load(msg.CorrelationId); ok {
				ch.(chan amqp.Delivery) <- msg
				log.Println("Message acknoledged from callback queue", msg.DeliveryTag)
			}
		}
	}()

	log.Println("listening to messages in gateway service")

	<-forever

	return nil
}
