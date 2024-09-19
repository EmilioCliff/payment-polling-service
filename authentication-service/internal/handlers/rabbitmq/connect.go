package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/EmilioCliff/payment-polling-app/authentication-service/internal/repository"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/pkg"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	AUTH_CONSUMER_NAME = "authentication_service"
)


type RabbitConn struct {
	conn   *amqp.Connection
	Config pkg.Config
	Maker pkg.JWTMaker

	UserRepository repository.UserRepository
}

func NewRabbitConn(config pkg.Config, tokenMaker pkg.JWTMaker) *RabbitConn {
	rabbit := &RabbitConn{
		Config: config,
		Maker: tokenMaker,
	}

	return rabbit
}

type Payload struct {
	Name string `json:"name"`
	Data any    `json:"data"`
}

func (r *RabbitConn) ConnectToRabbit() error {
	count := 0
	rollOff := time.Second
	var err error
	var connection *amqp.Connection
	for {
		connection, err = amqp.Dial(r.Config.RABBITMQ_URL)
		if err != nil {
			log.Println("failed to connect to rabbitmq", err)
			if count > 12 {
				return err
			}
			count++
			rollOff = time.Duration(math.Pow(float64(count), 2)) * time.Second
			time.Sleep(rollOff)
			continue
		}
		fmt.Println("Connected to rabbitmq")
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

	q, err := ch.QueueDeclare(
		r.Config.AUTH_QUEUE_NAME, // name
		false,                    // durable
		false,                    // delete when unused
		false,                    // exclusive
		false,                    // no-wait
		nil,                      // arguments
	)
	if err != nil {
		return err
	}

	for _, topic := range topics {
		if err := ch.QueueBind(
			q.Name,        // queue name
			topic,         // routing key
			r.Config.EXCH, // exchange
			false,
			nil,
		); err != nil {
			return err
		}
	}

	messages, err := ch.Consume(
		q.Name,             // queue
		AUTH_CONSUMER_NAME, // consumer
		false,              // auto ack
		false,              // exclusive
		false,              // no local
		false,              // no wait
		nil,                // args
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

			response := r.DistributeTask(payload)

			log.Printf("Message acknowledged from auth service: %v", msg.DeliveryTag)
			msg.Ack(false)

			count := 0
			for {
				err = ch.PublishWithContext(ctx,
					r.Config.EXCH, // exchange
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

func (r *RabbitConn) DistributeTask(payload Payload) []byte {
	switch payload.Name {
	case "register_user":
		dataBytes, err := json.Marshal(payload.Data)
		if err != nil {
			return r.errorRabbitMQResponse(pkg.Errorf(pkg.INTERNAL_ERROR, "failed to marshal request: %v", err))
		}

		var registerUserPayload RegisterUserRequest
		err = json.Unmarshal(dataBytes, &registerUserPayload)
		if err != nil {
			return r.errorRabbitMQResponse(pkg.Errorf(pkg.INTERNAL_ERROR, "failed to unmarshal request: %v", err))
		}

		return r.HandleRegisterUser(registerUserPayload)

	case "login_user":
		dataBytes, err := json.Marshal(payload.Data)
		if err != nil {
			return r.errorRabbitMQResponse(pkg.Errorf(pkg.INTERNAL_ERROR, "failed to marshal request: %v", err))
		}

		var loginUserPayload LoginUserRequest
		err = json.Unmarshal(dataBytes, &loginUserPayload)
		if err != nil {
			return r.errorRabbitMQResponse(pkg.Errorf(pkg.INTERNAL_ERROR, "failed to unmarshal request: %v", err))
		}

		return r.HandleLoginUser(loginUserPayload)

	default:
		// log unknow message
		return nil
	}
}
