package api

import (
	"context"
	"net/http"
	"time"

	pb "github.com/EmilioCliff/payment-polling-service/shared-grpc/pb"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"
)

func registerUserViagRPC(ctx *gin.Context, server *Server) {
	var req registerUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, server.errorResponse(err, "invalid request body"))
		return
	}

	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	rsp, err := server.authgRPClient.RegisterUser(c, &pb.RegisterUserRequest{
		Fullname:           req.FullName,
		Email:              req.Email,
		Password:           req.Password,
		PaydUsername:       req.PaydUsername,
		PaydPasswordApiKey: req.PaydPasswordApiKey,
		PaydUsernameApiKey: req.PaydUsernameApiKey,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			grpcCode := st.Code()
			grpcMessage := st.Message()
			httpCode, errRsp := server.grpcErrorResponse(grpcCode, grpcMessage)
			ctx.JSON(httpCode, errRsp)
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, server.errorResponse(err, "Non-gRPC error occurred"))
			return
		}
	}

	ctx.JSON(http.StatusOK, rsp)
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
		st, ok := status.FromError(err)
		if ok {
			grpcCode := st.Code()
			grpcMessage := st.Message()
			httpCode, errRsp := server.grpcErrorResponse(grpcCode, grpcMessage)
			ctx.JSON(httpCode, errRsp)
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, server.errorResponse(err, "Non-gRPC error occurred"))
			return
		}
	}

	ctx.JSON(http.StatusOK, rsp)
}
