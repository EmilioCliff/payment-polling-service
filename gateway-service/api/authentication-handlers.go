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

// registerUser sends a POST request to the authApp at "http://authApp:5000/auth/register"
// with the request body from the current context. It then reads the response body and
// unmarshals it into an authServiceRegisterResponse struct. If the status field of the
// response is false, it returns an error response. Otherwise, it returns the response
// body as a JSON object with a status code of 200.
//
// Parameters:
// - ctx: A gin Context object representing the current HTTP request and response.
//
// Return type: None.
func (server *Server) registerUser(ctx *gin.Context) {
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

// loginUser handles the login request from the client and sends a POST request to the authApp
// to authenticate the user. It reads the response body and unmarshals it into an
// authServiceLoginResponse struct. If the status field of the response is false, it returns
// an error response. Otherwise, it returns the response body as a JSON object with a status code
// of 200.
//
// Parameters:
// - ctx: A gin Context object representing the current HTTP request and response.
//
// Return type: None.
func (server *Server) loginUser(ctx *gin.Context) {
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
