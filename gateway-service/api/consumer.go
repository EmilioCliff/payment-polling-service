package api

import (
	"log"
)

const (
	INITIATE_PAYMENT = "gateway.initiate_payment"
)

func (server *Server) SetConsumer(topics []string) error {
	for _, topic := range topics {
		if err := server.amqpChannel.QueueBind(
			server.config.QUEUE_NAME, // queue name
			topic,                    // routing key
			server.config.EXCH,       // exchange
			false,
			nil,
		); err != nil {
			return err
		}
	}

	messages, err := server.amqpChannel.Consume(
		server.config.QUEUE_NAME, // queue
		"gateway_service",        // consumer
		true,                     // auto ack
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
		for msg := range messages {

			switch msg.CorrelationId {
			case "123":
				// receive the same message back that is sent
				log.Printf("this is the response(same as sent/testing: initiate payment) received from the message queue: %s, %s", msg.Body, msg.CorrelationId)
			}
		}
	}()

	log.Println("listening to messages in gateway service")

	<-forever

	return nil
}
