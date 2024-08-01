package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Payload struct {
	Name string `json:"name"`
	Data any    `json:"data"`
}

type initiatePaymentRequest struct {
	Amount        int64  `json:"amount"`
	PaymentMethod string `json:"payment_method"`
}

func (server *Server) initiatePayment(ctx *gin.Context) {
	var req initiatePaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, server.errorResponse(err, "invalid request body"))
		return
	}

	// payload, err := json.Marshal(req)
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, server.errorResponse(err, "error marshalling request body"))
	// 	return
	// }

	// data, err := ctx.GetRawData()
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, server.errorResponse(err, "error reading request body"))
	// 	return
	// }

	payload := Payload{
		Name: "initiate_payment",
		Data: req,
	}

	payloadRabitData, err := json.Marshal(payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, server.errorResponse(err, "error marshalling request body"))
		return
	}

	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// routing key is the same correlation id
	err = server.amqpChannel.PublishWithContext(c,
		server.config.EXCH,          // exchange
		"payments.initiate_payment", // routing key
		false,                       // mandatory
		false,                       // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: "123",
			ReplyTo:       INITIATE_PAYMENT,
			Body:          payloadRabitData,
		})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, server.errorResponse(err, "error communicating to the payment service via rabbitmq"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success-sent-raw-data": req})
}

func (server *Server) paymentStatus(ctx *gin.Context) {

}
