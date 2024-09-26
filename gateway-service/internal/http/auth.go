package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/EmilioCliff/payment-polling-app/gateway-service/internal/services"
	"github.com/EmilioCliff/payment-polling-app/gateway-service/pkg"
)

var _ services.HttpInterface = (*HTTPService)(nil)

type HTTPService struct {
	registerPath string
	loginPath    string
	config       pkg.Config
}

func NewHTTPService(config pkg.Config) *HTTPService {
	return &HTTPService{
		config:       config,
		registerPath: "auth/register",
		loginPath:    "auth/login",
	}
}

func (s *HTTPService) RegisterUserViaHttp(req services.RegisterUserRequest) (int, services.RegisterUserResponse) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return http.StatusInternalServerError, services.RegisterUserResponse{Message: "internal error", StatusCode: http.StatusInternalServerError}
	}

	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/%s", s.config.AUTH_HTTP_PORT, s.registerPath), bytes.NewBuffer(jsonData))
	if err != nil {
		return http.StatusInternalServerError, services.RegisterUserResponse{Message: "internal error", StatusCode: http.StatusInternalServerError}
	}

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return http.StatusInternalServerError, services.RegisterUserResponse{Message: "internal error", StatusCode: http.StatusInternalServerError}
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return http.StatusInternalServerError, services.RegisterUserResponse{Message: "internal error", StatusCode: http.StatusInternalServerError}
	}

	var authServiceResponse services.RegisterUserResponse

	err = json.Unmarshal(responseBody, &authServiceResponse)
	if err != nil {
		return http.StatusInternalServerError, services.RegisterUserResponse{Message: "internal error", StatusCode: http.StatusInternalServerError}
	}

	if response.StatusCode != http.StatusOK {
		return response.StatusCode, services.RegisterUserResponse{Message: authServiceResponse.Message, StatusCode: authServiceResponse.StatusCode}
	}

	return http.StatusOK, authServiceResponse
}

func (s *HTTPService) LoginUserViaHttp(req services.LoginUserRequest) (int, services.LoginUserResponse) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return http.StatusInternalServerError, services.LoginUserResponse{Message: "internal error", StatusCode: http.StatusInternalServerError}
	}

	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/%s", s.config.AUTH_HTTP_PORT, s.loginPath), bytes.NewBuffer(jsonData))
	if err != nil {
		return http.StatusInternalServerError, services.LoginUserResponse{Message: "internal error", StatusCode: http.StatusInternalServerError}
	}

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return http.StatusInternalServerError, services.LoginUserResponse{Message: "internal error", StatusCode: http.StatusInternalServerError}
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return http.StatusInternalServerError, services.LoginUserResponse{Message: "internal error", StatusCode: http.StatusInternalServerError}
	}

	var jsonFromAuthService services.LoginUserResponse

	err = json.Unmarshal(responseBody, &jsonFromAuthService)
	if err != nil {
		return http.StatusInternalServerError, services.LoginUserResponse{Message: "internal error", StatusCode: http.StatusInternalServerError}
	}

	if response.StatusCode != http.StatusOK {
		return response.StatusCode, services.LoginUserResponse{Message: jsonFromAuthService.Message, StatusCode: jsonFromAuthService.StatusCode}
	}

	return http.StatusOK, jsonFromAuthService
}
