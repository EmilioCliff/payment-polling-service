package rabbitmq

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/EmilioCliff/payment-polling-app/gateway-service/internal/services"
	"github.com/EmilioCliff/payment-polling-app/gateway-service/pkg"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (r *RabbitHandler) RegisterUserViaRabbit(req services.RegisterUserRequest) (int, gin.H) {
	payload := services.Payload{
		Name: "register_user",
		Data: req,
	}

	payloadRabitData, err := json.Marshal(payload)
	if err != nil {
		return http.StatusInternalServerError, pkg.ErrorResponse(err, "error marshalling request body")
	}

	correlationID := uuid.New().String()
	responseChannel := make(chan amqp.Delivery, 1)
	defer close(responseChannel)

	r.RspMap.Set(correlationID, responseChannel)
	defer r.RspMap.Delete(correlationID)

	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = r.Channel.PublishWithContext(c,
		r.config.EXCH,                  // exchange
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
		return http.StatusInternalServerError, pkg.ErrorResponse(err, "error communicating to the authentication service via rabbitmq")
	}

	select {
	case msg := <-responseChannel:
		if msg.CorrelationId == correlationID {
			var authResp services.RegisterUserResponse

			err := json.Unmarshal(msg.Body, &authResp)
			if err != nil {
				return http.StatusInternalServerError, gin.H{"error": "Failed to parse response"}
			}

			if authResp.Message != "" {
				return authResp.StatusCode, gin.H{"error": authResp.Message}
			}

			return http.StatusOK, gin.H{"response": authResp}
		}
	case <-time.After(5 * time.Second):
		return http.StatusRequestTimeout, gin.H{"error": "Timeout waiting for response"}
	}

	return http.StatusInternalServerError, gin.H{"error": "unknown error"}
}
func (r *RabbitHandler) LoginUserViaRabbit(req services.LoginUserRequest) (int, gin.H) {
	payload := services.Payload{
		Name: "login_user",
		Data: req,
	}

	payloadRabitData, err := json.Marshal(payload)
	if err != nil {
		return http.StatusInternalServerError, pkg.ErrorResponse(err, "error marshalling request body")
	}

	correlationID := uuid.New().String()
	responseChannel := make(chan amqp.Delivery, 1)
	defer close(responseChannel)

	r.RspMap.Set(correlationID, responseChannel)
	defer r.RspMap.Delete(correlationID)

	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = r.Channel.PublishWithContext(c,
		r.config.EXCH,               // exchange
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
		return http.StatusInternalServerError, pkg.ErrorResponse(err, "error communicating to the authentication service via rabbitmq")
	}

	select {
	case msg := <-responseChannel:
		if msg.CorrelationId == correlationID {
			var loginResp services.LoginUserResponse

			err := json.Unmarshal(msg.Body, &loginResp)
			if err != nil {
				return http.StatusInternalServerError, gin.H{"error": "Failed to parse response"}
			}

			if loginResp.Message != "" {
				return loginResp.StatusCode, gin.H{"error": loginResp.Message}
			}

			return http.StatusOK, gin.H{"response": loginResp}
		}

	case <-time.After(5 * time.Second):
		return http.StatusRequestTimeout, gin.H{"error": "Timeout waiting for response"}
	}

	return http.StatusInternalServerError, gin.H{"error": "unknown error"}
}
