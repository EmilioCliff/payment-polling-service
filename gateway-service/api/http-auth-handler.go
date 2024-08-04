package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type authRegisterResponse struct {
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type authServiceRegisterResponse struct {
	Status  bool                 `json:"status"`
	Message string               `json:"message,omitempty"`
	Data    authRegisterResponse `json:"data,omitempty"`
}

func registerUserViaHTTP(ctx *gin.Context, server *Server) {
	request, err := http.NewRequest("POST", "http://authApp:5000/auth/register", ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, server.errorResponse(err, "Couldn't create request to authApp"))
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, server.errorResponse(err, "Couldn't send request to authApp"))
		return
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, server.errorResponse(err, "Couldn't read response body from authApp"))
		return
	}

	var authServiceResponse authServiceRegisterResponse

	err = json.Unmarshal(responseBody, &authServiceResponse)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, server.errorResponse(err, "Couldn't unmarshal response body from authApp"))
		return
	}

	if !authServiceResponse.Status {
		ctx.JSON(response.StatusCode, server.errorResponse(errors.New(authServiceResponse.Message), "Failed to register new user"))
		return
	}

	ctx.JSON(http.StatusOK, authServiceResponse)
}

type authLoginResponse struct {
	AccessToken string    `json:"access_token"`
	FullName    string    `json:"full_name"`
	Email       string    `json:"email"`
	CreatedAt   time.Time `json:"created_at"`
}

type authServiceLoginResponse struct {
	Status  bool              `json:"status"`
	Message string            `json:"message,omitempty"`
	Data    authLoginResponse `json:"data,omitempty"`
}

func loginUserViaHTTP(ctx *gin.Context, server *Server) {
	request, err := http.NewRequest("POST", "http://authApp:5000/auth/login", ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, server.errorResponse(err, "Couldn't create request to authApp"))
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, server.errorResponse(err, "Couldn't send request to authApp"))
		return
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, server.errorResponse(err, "Couldn't read response body from authApp"))
		return
	}

	var jsonFromAuthService authServiceLoginResponse

	err = json.Unmarshal(responseBody, &jsonFromAuthService)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, server.errorResponse(err, "Couldn't unmarshal response body from authApp"))
		return
	}

	if !jsonFromAuthService.Status {
		ctx.JSON(response.StatusCode, server.errorResponse(errors.New(jsonFromAuthService.Message), "Failed to register new user"))
		return
	}

	ctx.JSON(http.StatusOK, jsonFromAuthService)
}
