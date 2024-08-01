package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/EmilioCliff/payment-polling-app/gateway-service/pb"
	"github.com/gin-gonic/gin"
)

func (server *Server) registerUser(ctx *gin.Context) {
	registerUserViagRPC(ctx, server)
	// registerUserViaHTTP(ctx, server)
}

func (server *Server) loginUser(ctx *gin.Context) {
	loginUserViagRPC(ctx, server)
	// loginUserViaHTTP(ctx, server)
}

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

type registerUserRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func registerUserViagRPC(ctx *gin.Context, server *Server) {
	var req registerUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, server.errorResponse(err, "invalid request body"))
		return
	}

	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	rsp, err := server.authgRPClient.RegisterUser(c, &pb.RegisterUserRequest{
		Fullname: req.FullName,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, server.errorResponse(err, "error from gRPC server"))
		return
	}

	ctx.JSON(http.StatusOK, rsp)
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

type loginUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func loginUserViagRPC(ctx *gin.Context, server *Server) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, server.errorResponse(err, "invalid request body"))
		return
	}

	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	rsp, err := server.authgRPClient.LoginUser(c, &pb.LoginUserRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, server.errorResponse(err, "error from gRPC server"))
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}
