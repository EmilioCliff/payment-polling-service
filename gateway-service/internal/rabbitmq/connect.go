package rabbitmq

import (
	"log"
	"sync"

	"github.com/EmilioCliff/payment-polling-app/gateway-service/internal/services"
	"github.com/EmilioCliff/payment-polling-app/gateway-service/pkg"
	amqp "github.com/rabbitmq/amqp091-go"
)

var _ services.RabbitInterface = (*RabbitHandler)(nil)

type RabbitHandler struct {
	Channel *amqp.Channel
	RspMap  *responseMap
	config  pkg.Config
}

type responseMap struct {
	mu   sync.RWMutex
	data map[string]chan amqp.Delivery
}

func NewRabbitHandler(channel *amqp.Channel, config pkg.Config) *RabbitHandler {
	return &RabbitHandler{
		RspMap:  NewResponseMap(),
		config:  config,
		Channel: channel,
	}
}

func NewResponseMap() *responseMap {
	return &responseMap{
		data: make(map[string]chan amqp.Delivery),
	}
}

func (r *RabbitHandler) SetConsumer(topics []string) error {
	q, err := r.Channel.QueueDeclare(
		r.config.EXCLUSIVE_QUEUE_NAME, // name
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
		if err := r.Channel.QueueBind(
			q.Name,        // queue name
			topic,         // routing key
			r.config.EXCH, // exchange
			false,
			nil,
		); err != nil {
			return err
		}
	}

	messages, err := r.Channel.Consume(
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
			if ch, ok := r.RspMap.Get(msg.CorrelationId); ok {
				ch <- msg
				log.Println("Message acknoledged from callback queue", msg.DeliveryTag)
				// log.Println(msg)
			}
		}
	}()

	log.Println("listening to messages in gateway service")

	<-forever

	return nil
}

func (rm *responseMap) Set(correlationID string, channel chan amqp.Delivery) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.data[correlationID] = channel
}

func (rm *responseMap) Get(correlationID string) (chan amqp.Delivery, bool) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	ch, exists := rm.data[correlationID]
	return ch, exists
}

func (rm *responseMap) Delete(correlationID string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	delete(rm.data, correlationID)
}
