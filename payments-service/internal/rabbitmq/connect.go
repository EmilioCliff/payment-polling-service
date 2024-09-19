package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/EmilioCliff/payment-polling-app/payment-service/internal/postgres"
	"github.com/EmilioCliff/payment-polling-app/payment-service/internal/workers"
	"github.com/EmilioCliff/payment-polling-app/payment-service/pkg"
	"github.com/EmilioCliff/payment-polling-service/shared-grpc/pb"
	"github.com/hibiken/asynq"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	PAYMENT_CONSUMER_NAME = "payment_service"
)

type RabbitConn struct {
	distributor workers.TaskDistributor
	conn        *amqp.Connection
	config      pkg.Config
	client      pb.AuthenticationServiceClient
	store       postgres.Store
}

func NewRabbitConn() (*RabbitConn, error) {
	rabbit := &RabbitConn{}

	config, err := pkg.LoadConfig(".")
	if err != nil {
		return nil, err
	}

	distributor := workers.NewRedisTaskDistributor(asynq.RedisClientOpt{
		Addr: config.REDDIS_ADDR,
		DB:   1,
	})

	conn, err := rabbit.connectToRabbit(config.RABBITMQ_URL)
	if err != nil {
		return nil, err
	}

	gRPCconn, err := grpc.NewClient(config.AUTH_GRPC_URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	store, err := postgres.GetStore()
	if err != nil {
		return nil, err
	}

	rabbit.config = config
	rabbit.distributor = distributor
	rabbit.conn = conn
	rabbit.client = pb.NewAuthenticationServiceClient(gRPCconn)
	rabbit.store = store

	return rabbit, nil
}

type Payload struct {
	Name string `json:"name"`
	Data any    `json:"data"`
}

func (r *RabbitConn) connectToRabbit(rabitURL string) (*amqp.Connection, error) {
	count := 0
	rollOff := 1 * time.Second
	var err error
	var connection *amqp.Connection
	for {
		connection, err = amqp.Dial(rabitURL)
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

	messages, err := ch.Consume(
		q.Name,                // queue
		PAYMENT_CONSUMER_NAME, // consumer
		false,                 // auto ack
		false,                 // exclusive
		false,                 // no local
		false,                 // no wait
		nil,                   // args
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
				log.Printf("failed to unmarshal message in payment service: %s", err)
				msg.Nack(false, true)
				return
			}

			response := r.distributeTask(payload)

			log.Printf("Message acknowledged from payment service: %v", msg.DeliveryTag)
			msg.Ack(false)

			count := 0
			for {
				err = ch.PublishWithContext(ctx,
					r.config.EXCH, // exchange
					msg.ReplyTo,   // routing key
					false,         // mandatory
					false,         // immediate
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
						// log to failed to send response
						log.Printf("failed to send response: %s", err)
						return
					}
				}
			}
		}
	}()

	log.Println("listening to messages in authentication service")

	<-forever

	return nil
}

func (r *RabbitConn) distributeTask(payload Payload) []byte {
	switch payload.Name {
	case "initiate_payment":
		dataBytes, err := json.Marshal(payload.Data)
		if err != nil {
			return r.errorRabbitMQResponse(pkg.Errorf(pkg.INTERNAL_ERROR, "failed to marshal request: %v", err))
		}

		var initiatePaymentPayload initiatePaymentRequest
		err = json.Unmarshal(dataBytes, &initiatePaymentPayload)
		if err != nil {
			return r.errorRabbitMQResponse(pkg.Errorf(pkg.INTERNAL_ERROR, "failed to unmarshal request: %v", err))
		}

		return r.handleInitiatePayment(initiatePaymentPayload)

	case "polling_transaction":
		dataBytes, err := json.Marshal(payload.Data)
		if err != nil {
			return r.errorRabbitMQResponse(pkg.Errorf(pkg.INTERNAL_ERROR, "failed to marshal request: %v", err))
		}

		var pollingTransactionPayload postgres.PollingTransactionRequest
		err = json.Unmarshal(dataBytes, &pollingTransactionPayload)
		if err != nil {
			return r.errorRabbitMQResponse(pkg.Errorf(pkg.INTERNAL_ERROR, "failed to unmarshal request: %v", err))
		}

		return r.handlePollingTransaction(pollingTransactionPayload)

	default:
		// log unknow message
		return nil
	}
}
