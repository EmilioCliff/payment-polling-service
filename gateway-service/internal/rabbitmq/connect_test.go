package rabbitmq

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/EmilioCliff/payment-polling-app/gateway-service/pkg"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/require"
	mq "github.com/testcontainers/testcontainers-go/modules/rabbitmq"
)

type TestRabbitHandler struct {
	rabbit    *RabbitHandler
	container *mq.RabbitMQContainer
	ctx       context.Context
}

func NewTestRabbitHandler() (*TestRabbitHandler, error) {
	ctx := context.Background()

	rabbitmqContainer, err := mq.Run(ctx,
		"rabbitmq:3.12.11-management-alpine",
	)
	if err != nil {
		log.Printf("failed to start container: %s", err)

		return nil, err
	}

	log.Println(rabbitmqContainer.HttpURL(ctx))

	connString, err := rabbitmqContainer.AmqpURL(ctx)
	if err != nil {
		return nil, err
	}

	conn, err := ConnectToRabit(connString)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		"events", // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return nil, err
	}

	r := NewRabbitService(ch, pkg.Config{
		EXCLUSIVE_QUEUE_NAME:  "gateway_queue",
		EXCH:                  "events",
		GATEWAY_CONSUMER_NAME: "gateway_service",
	})

	return &TestRabbitHandler{
		ctx:       ctx,
		container: rabbitmqContainer,
		rabbit:    r,
	}, nil
}

func TestRabbitHandler_TestSetConsumer(t *testing.T) {
	pkg.SkipCI(t)

	testRabbit, err := NewTestRabbitHandler()
	require.NoError(t, err)

	defer func() {
		// close the channel and terminate the container
		testRabbit.rabbit.Channel.Close()

		if err := testRabbit.container.Terminate(testRabbit.ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()

	errCh := make(chan error, 1)

	readyChan := make(chan struct{}, 1)

	go func(errCh chan error, readyChan chan struct{}) {
		err = testRabbit.rabbit.SetConsumer(
			[]string{
				"gateway.initiate_payment",
				"gateway.poll_payments",
				"gateway.register_user",
				"gateway.login_user",
			}, readyChan,
		)
		errCh <- err
	}(errCh, readyChan)

	<-readyChan
	close(readyChan)

	corrID := "test-correlation-id"
	ch := testRabbit.rabbit.Channel
	err = ch.Publish(
		"events",                   // exchange
		"gateway.initiate_payment", // routing key
		false,                      // mandatory
		false,                      // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			Body:          []byte("test message"),
			CorrelationId: corrID,
		},
	)
	require.NoError(t, err)

	msgCh := make(chan amqp.Delivery, 1)
	testRabbit.rabbit.RspMap.Set(corrID, msgCh)

	select {
	case msg := <-msgCh:
		require.Equal(t, "test message", string(msg.Body))
		close(testRabbit.rabbit.forever)
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for message to be consumed")
	}

	select {
	case err := <-errCh:
		require.NoError(t, err)
		log.Println(err)
		close(errCh)

		return
	case <-time.After(5 * time.Second):
		t.Fatal("SetConsumer timed out")
	}
}
