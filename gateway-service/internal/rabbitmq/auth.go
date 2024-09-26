package rabbitmq

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/EmilioCliff/payment-polling-app/gateway-service/internal/services"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (r *RabbitHandler) RegisterUserViaRabbit(req services.RegisterUserRequest) (int, services.RegisterUserResponse) {
	payload := services.Payload{
		Name: "register_user",
		Data: req,
	}

	payloadRabitData, err := json.Marshal(payload)
	if err != nil {
		return http.StatusInternalServerError, services.RegisterUserResponse{Message: "internal error", StatusCode: http.StatusInternalServerError}
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
		return http.StatusInternalServerError, services.RegisterUserResponse{Message: "internal error", StatusCode: http.StatusInternalServerError}
	}

	select {
	case msg := <-responseChannel:
		if msg.CorrelationId == correlationID {
			var authResp services.RegisterUserResponse

			err := json.Unmarshal(msg.Body, &authResp)
			if err != nil {
				return http.StatusInternalServerError, services.RegisterUserResponse{
					Message:    "internal error",
					StatusCode: http.StatusInternalServerError,
				}
			}

			if authResp.Message != "" {
				return authResp.StatusCode, services.RegisterUserResponse{Message: authResp.Message, StatusCode: authResp.StatusCode}
			}

			return http.StatusOK, authResp
		}
	case <-time.After(5 * time.Second):
		return http.StatusRequestTimeout, services.RegisterUserResponse{
			Message:    "timeout waiting for response. Try again",
			StatusCode: http.StatusInternalServerError,
		}
	}

	return http.StatusInternalServerError, services.RegisterUserResponse{Message: "internal error", StatusCode: http.StatusInternalServerError}
}

func (r *RabbitHandler) LoginUserViaRabbit(req services.LoginUserRequest) (int, services.LoginUserResponse) {
	payload := services.Payload{
		Name: "login_user",
		Data: req,
	}

	payloadRabitData, err := json.Marshal(payload)
	if err != nil {
		return http.StatusInternalServerError, services.LoginUserResponse{Message: "internal error", StatusCode: http.StatusInternalServerError}
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
		return http.StatusInternalServerError, services.LoginUserResponse{Message: "internal error", StatusCode: http.StatusInternalServerError}
	}

	select {
	case msg := <-responseChannel:
		if msg.CorrelationId == correlationID {
			var loginResp services.LoginUserResponse

			err := json.Unmarshal(msg.Body, &loginResp)
			if err != nil {
				return http.StatusInternalServerError, services.LoginUserResponse{
					Message:    "internal error",
					StatusCode: http.StatusInternalServerError,
				}
			}

			if loginResp.Message != "" {
				return loginResp.StatusCode, services.LoginUserResponse{Message: loginResp.Message, StatusCode: loginResp.StatusCode}
			}

			return http.StatusOK, loginResp
		}

	case <-time.After(5 * time.Second):
		return http.StatusRequestTimeout, services.LoginUserResponse{
			Message:    "timeout waiting for response. Try again",
			StatusCode: http.StatusInternalServerError,
		}
	}

	return http.StatusInternalServerError, services.LoginUserResponse{Message: "internal error", StatusCode: http.StatusInternalServerError}
}
