package rabbitmq

import (
	"encoding/json"
	"log"

	"github.com/EmilioCliff/payment-polling-app/authentication-service/pkg"
)

func (r *RabbitConn) errorRabbitMQResponse(pkgErr *pkg.Error) []byte {
	jsonResponse, err := json.Marshal(pkgErr)
	if err != nil {
		log.Printf("Failed to marshal error response: %v", err)
		return []byte(`{"status": false}`)
	}

	return jsonResponse
}
