package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

// "authentication.register_user", "authentication.login_user"

func (server *Server) registerUserViaRabbit(ctx *gin.Context) {
	log.Println("in register handler for rabbit")
	var req registerUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println("failed to bind json")
		ctx.JSON(http.StatusBadRequest, server.errorResponse(err, "invalid request body"))
		return
	}

	payload := Payload{
		Name: "register_user",
		Data: req,
	}

	payloadRabitData, err := json.Marshal(payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, server.errorResponse(err, "error marshalling request body"))
		return
	}

	correlationID := uuid.New().String()
	responseChannel := make(chan amqp.Delivery, 1)
	defer close(responseChannel)

	server.responseMap.Store(correlationID, responseChannel)
	defer server.responseMap.Delete(correlationID)

	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("payload from register", payload)

	err = server.amqpChannel.PublishWithContext(c,
		server.config.EXCH,             // exchange
		"authentication.register_user", // routing key
		false,                          // mandatory
		false,                          // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: correlationID,
			ReplyTo:       "gateway.register_user",
			Body:          payloadRabitData,
		})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, server.errorResponse(err, "error communicating to the authentication service via rabbitmq"))
		return
	}

	select {
	case msg := <-responseChannel:
		if msg.CorrelationId == correlationID {
			ctx.JSON(http.StatusOK, gin.H{"auth response": string(msg.Body)})
		}
	case <-time.After(5 * time.Second):
		ctx.JSON(http.StatusRequestTimeout, gin.H{"error": "Timeout waiting for response"})
	}
}

func loginUserViaRabbit(ctx *gin.Context, server *Server) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, server.errorResponse(err, "invalid request body"))
		return
	}

	payload := Payload{
		Name: "login_user",
		Data: req,
	}

	payloadRabitData, err := json.Marshal(payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, server.errorResponse(err, "error marshalling request body"))
		return
	}

	correlationID := uuid.New().String()
	responseChannel := make(chan amqp.Delivery, 1)
	defer close(responseChannel)

	server.responseMap.Store(correlationID, responseChannel)
	defer server.responseMap.Delete(correlationID)

	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.amqpChannel.PublishWithContext(c,
		server.config.EXCH,          // exchange
		"authentication.login_user", // routing key
		false,                       // mandatory
		false,                       // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: correlationID,
			ReplyTo:       "gateway.login_user",
			Body:          payloadRabitData,
		})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, server.errorResponse(err, "error communicating to the authentication service via rabbitmq"))
		return
	}

	select {
	case msg := <-responseChannel:
		if msg.CorrelationId == correlationID {
			ctx.JSON(http.StatusOK, gin.H{"response": string(msg.Body)})
		}
	case <-time.After(5 * time.Second):
		ctx.JSON(http.StatusRequestTimeout, gin.H{"error": "Timeout waiting for response"})
	}
}
