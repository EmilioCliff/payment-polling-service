package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/EmilioCliff/payment-polling-app/gateway-service/internal/services"
	"github.com/EmilioCliff/payment-polling-app/gateway-service/pkg"
	"github.com/gin-gonic/gin"
)

func (s *HttpServer) RegisterUserViaHttp(req services.RegisterUserRequest) (int, gin.H) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return http.StatusInternalServerError, pkg.ErrorResponse(err, "Couldn't marshal request data")
	}

	request, err := http.NewRequest("POST", "http://authApp:5000/auth/register", bytes.NewBuffer(jsonData))
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
		return http.StatusInternalServerError, pkg.ErrorResponse(err, "Couldn't read response body from authApp")
	}

	var authServiceResponse services.RegisterUserResponse

	err = json.Unmarshal(responseBody, &authServiceResponse)
	if err != nil {
		return http.StatusInternalServerError, pkg.ErrorResponse(err, "Couldn't unmarshal response body from authApp")
	}

	if response.StatusCode != http.StatusOK {
		return response.StatusCode, pkg.ErrorResponse(errors.New(authServiceResponse.Message), "Failed to register new user")
	}

	return http.StatusOK, gin.H{"response": authServiceResponse}
}

func (s *HttpServer) LoginUserViaHttp(req services.LoginUserRequest) (int, gin.H) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return http.StatusInternalServerError, pkg.ErrorResponse(err, "Couldn't marshal request data")
	}

	request, err := http.NewRequest("POST", "http://authApp:5000/auth/login", bytes.NewBuffer(jsonData))
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
		return http.StatusInternalServerError, pkg.ErrorResponse(err, "Couldn't read response body from authApp")
	}

	var jsonFromAuthService services.LoginUserResponse

	err = json.Unmarshal(responseBody, &jsonFromAuthService)
	if err != nil {
		return http.StatusInternalServerError, pkg.ErrorResponse(err, "Couldn't unmarshal response body from authApp")
	}

	if response.StatusCode != http.StatusOK {
		return response.StatusCode, pkg.ErrorResponse(errors.New(jsonFromAuthService.Message), "Failed to login user")
	}

	return http.StatusOK, gin.H{"response": jsonFromAuthService}
}
