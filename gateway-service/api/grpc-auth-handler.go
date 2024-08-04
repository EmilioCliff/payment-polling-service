package api

import (
	"context"
	"net/http"
	"time"

	"github.com/EmilioCliff/payment-polling-app/gateway-service/pb"
	"github.com/gin-gonic/gin"
)

type registerUserRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
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
