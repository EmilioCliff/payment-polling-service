package rabbitmq

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/EmilioCliff/payment-polling-app/payment-service/pkg"
)

type errorResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func (r *RabbitConn) errorRabbitMQResponse(pkgErr *pkg.Error) []byte {
	errorRsp := errorResponse{
		StatusCode: r.convertPkgError(pkgErr),
		Message:    pkgErr.Message,
	}

	jsonResponse, err := json.Marshal(errorRsp)
	if err != nil {
		log.Printf("Failed to marshal error response: %v", err)
		return []byte(`{"status": false}`)
	}

	return jsonResponse
}

func (r RabbitConn) convertPkgError(err *pkg.Error) int {
	switch err.Code {
	case pkg.ALREADY_EXISTS_ERROR:
		return http.StatusConflict
	case pkg.INTERNAL_ERROR:
		return http.StatusInternalServerError
	case pkg.INVALID_ERROR:
		return http.StatusBadRequest
	case pkg.NOT_FOUND_ERROR:
		return http.StatusNotFound
	case pkg.NOT_IMPLEMENTED_ERROR:
		return http.StatusNotImplemented
	case pkg.AUTHENTICATION_ERROR:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
