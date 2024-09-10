package http

import (
	"net/http"

	"github.com/EmilioCliff/payment-polling-app/authentication-service/internal/postgres"
	"github.com/gin-gonic/gin"
)

func (s *HttpServer) handleRegisterUser(ctx *gin.Context) {
	var req postgres.RegisterUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rsp, err := s.store.RegisterUser(ctx, req)
	if err != nil {
		ctx.JSON(s.convertPkgError(err), gin.H{"status_code": s.convertPkgError(err), "message": err.Message})
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (s *HttpServer) handleLoginUser(ctx *gin.Context) {
	var req postgres.LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rsp, err := s.store.LoginUser(ctx, req)
	if err != nil {
		ctx.JSON(s.convertPkgError(err), gin.H{"status_code": s.convertPkgError(err), "message": err.Message})
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}
