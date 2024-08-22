package app

import (
	db "github.com/EmilioCliff/payment-polling-app/payment-service/db/sqlc"
	"github.com/EmilioCliff/payment-polling-app/payment-service/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

type App struct {
	conn   *amqp.Connection
	config utils.Config
	store  *db.Queries
}

func NewApp(conn *amqp.Connection, config utils.Config, store *db.Queries) *App {
	app := App{
		config: config,
		conn:   conn,
		store:  store,
	}

	return &app
}
