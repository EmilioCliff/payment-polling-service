package app

import (
	"github.com/EmilioCliff/payment-polling-app/payment-service/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

type App struct {
	conn   *amqp.Connection
	config utils.Config
}

func NewApp(conn *amqp.Connection, config utils.Config) *App {
	app := App{
		config: config,
		conn:   conn,
	}

	return &app
}
