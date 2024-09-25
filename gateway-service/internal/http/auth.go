package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/EmilioCliff/payment-polling-app/gateway-service/internal/services"
	"github.com/EmilioCliff/payment-polling-app/gateway-service/pkg"
	"github.com/gin-gonic/gin"
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

func (s *HTTPService) RegisterUserViaHttp(req services.RegisterUserRequest) (int, gin.H) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return http.StatusInternalServerError, pkg.ErrorResponse(err, "Couldn't marshal request data")
	}

	request, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("http://%s/%s", s.config.AUTH_HTTP_PORT, s.registerPath),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return http.StatusInternalServerError, pkg.ErrorResponse(err, "Couldn't create request to authApp")
	}

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return http.StatusInternalServerError, pkg.ErrorResponse(err, "Couldn't send request to authApp")
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return http.StatusInternalServerError, pkg.ErrorResponse(
			err,
			"Couldn't read response body from authApp",
		)
	}

	var authServiceResponse services.RegisterUserResponse

	err = json.Unmarshal(responseBody, &authServiceResponse)
	if err != nil {
		return http.StatusInternalServerError, pkg.ErrorResponse(
			err,
			"Couldn't unmarshal response body from authApp",
		)
	}

	if response.StatusCode != http.StatusOK {
		return response.StatusCode, pkg.ErrorResponse(
			errors.New(authServiceResponse.Message),
			"Failed to register new user",
		)
	}

	return http.StatusOK, gin.H{"response": authServiceResponse}
}

func (s *HTTPService) LoginUserViaHttp(req services.LoginUserRequest) (int, gin.H) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return http.StatusInternalServerError, pkg.ErrorResponse(err, "Couldn't marshal request data")
	}

	request, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("http://%s/%s", s.config.AUTH_HTTP_PORT, s.loginPath),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return http.StatusInternalServerError, pkg.ErrorResponse(err, "Couldn't create request to authApp")
	}

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return http.StatusInternalServerError, pkg.ErrorResponse(err, "Couldn't send request to authApp")
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return http.StatusInternalServerError, pkg.ErrorResponse(
			err,
			"Couldn't read response body from authApp",
		)
	}

	var jsonFromAuthService services.LoginUserResponse

	err = json.Unmarshal(responseBody, &jsonFromAuthService)
	if err != nil {
		return http.StatusInternalServerError, pkg.ErrorResponse(
			err,
			"Couldn't unmarshal response body from authApp",
		)
	}

	if response.StatusCode != http.StatusOK {
		return response.StatusCode, pkg.ErrorResponse(
			errors.New(jsonFromAuthService.Message),
			"Failed to login user",
		)
	}

	return http.StatusOK, gin.H{"response": jsonFromAuthService}
}
