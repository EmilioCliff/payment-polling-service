package rabbitmq

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"time"

	"github.com/EmilioCliff/payment-polling-app/payment-service/internal/repository"
	"github.com/EmilioCliff/payment-polling-app/payment-service/internal/workers"
	"github.com/EmilioCliff/payment-polling-app/payment-service/pkg"
	"github.com/EmilioCliff/payment-polling-service/shared-grpc/pb"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Payload struct {
	Name string `json:"name"`
	Data []byte `json:"data"`
}

type RabbitConn struct {
	conn   *amqp.Connection
	config pkg.Config

	client                pb.AuthenticationServiceClient
	Distributor           workers.TaskDistributor
	TransactionRepository repository.TransactionRepository
}

func NewRabbitConn(config pkg.Config, client pb.AuthenticationServiceClient) *RabbitConn {
	return &RabbitConn{
		config: config,
		client: client,
	}
}

func (r *RabbitConn) ConnectToRabbit() error {
	count := 0
	maxRetries := 12

	var err error

	var connection *amqp.Connection

	for {
		connection, err = amqp.Dial(r.config.RABBITMQ_URL)
		if err != nil {
			log.Println("failed to connect to rabbitmq", err)

			if count > maxRetries {
				return err
			}

			count++
			rollOff := time.Duration(math.Pow(float64(count), 2)) * time.Second
			time.Sleep(rollOff)

			continue
		}

		log.Println("Connected to rabbitmq")

		break
	}

	r.conn = connection

	return nil
}

func (r *RabbitConn) SetConsumer(topics []string) error {
	ch, err := r.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		r.config.EXCH, // name
		"topic",       // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(
		r.config.PAYMENT_QUEUE_NAME, // name
		false,                       // durable
		false,                       // delete when unused
		false,                       // exclusive
		false,                       // no-wait
		nil,                         // arguments
	)
	if err != nil {
		return err
	}

	for _, topic := range topics {
		if err := ch.QueueBind(
			q.Name,        // queue name
			topic,         // routing key
			r.config.EXCH, // exchange
			false,
			nil,
		); err != nil {
			return err
		}
	}

	msgs, err := ch.Consume(
		q.Name,                         // queue
		r.config.PAYMENT_CONSUMER_NAME, // consumer
		false,                          // auto-ack
		false,                          // exclusive
		false,                          // no-local
		false,                          // no-wait
		nil,                            // args
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		for d := range msgs {
			var payload Payload

			err := json.Unmarshal(d.Body, &payload)
			if err != nil {
				log.Printf("failed to unmarshal message in auth rabbit: %s", err)

				_ = d.Nack(false, true)

				return
			}

			response := r.distributeTask(payload)

			log.Printf("Message acknowledged from payment service: %v", d.DeliveryTag)

			_ = d.Ack(false)

			count := 0
			maxRetries := 5

			for {
				err = ch.PublishWithContext(ctx,
					r.config.EXCH, // exchange
					d.ReplyTo,     // routing key
					false,         // mandatory
					false,         // immediate
					amqp.Publishing{
						ContentType:   "text/plain",
						CorrelationId: d.CorrelationId,
						Body:          response,
					},
				)
				if err == nil {
					break
				}

				count++

				if count > maxRetries {
					// log to failed to send response
					log.Printf("failed to send response: %s", err)

					return
				}
			}
		}
	}()

	log.Println("listening to messages in payment service")

	<-forever

	return nil
}

func (r *RabbitConn) distributeTask(payload Payload) []byte {
	switch payload.Name {
	case "initiate_payment":
		var initiatePaymentPayload initiatePaymentRequest

		err := json.Unmarshal(payload.Data, &initiatePaymentPayload)
		if err != nil {
			return r.errorRabbitMQResponse(pkg.Errorf(pkg.INTERNAL_ERROR, "%v", err))
		}

		return r.handleInitiatePayment(initiatePaymentPayload)

	case "polling_transaction":
		var pollingTransactionPayload pollingTransactionRequest

		err := json.Unmarshal(payload.Data, &pollingTransactionPayload)
		if err != nil {
			return r.errorRabbitMQResponse(pkg.Errorf(pkg.INTERNAL_ERROR, "%v", err))
		}

		return r.handlePollingTransaction(pollingTransactionPayload)

	default:
		// log unknow message
		return nil
	}
}
