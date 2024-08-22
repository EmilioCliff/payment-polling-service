package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type registerUserRequest struct {
	FullName       string `json:"full_name" binding:"required"`
	PaydUsername   string `json:"payd_username" binding:"required"`
	Email          string `json:"email" binding:"required"`
	Password       string `json:"password" binding:"required"`
	UsernameApiKey string `json:"payd_username_api_key" binding:"required"`
	PasswordApiKey string `json:"payd_password_api_key" binding:"required"`
}

func (server *Server) registerUser(ctx *gin.Context) {
	var req registerUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, server.errorResponse(err, "bad request to authApp"))
		return
	}

	rsp, errRsp := server.registerUserGeneral(req, ctx)
	if !errRsp.Status {
		ctx.JSON(errRsp.HttpStatus, server.errorResponse(errRsp.Error, errRsp.Message))
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}

type loginUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, server.errorResponse(err, "bad request to authApp"))
		return
	}

	rsp, errRsp := server.loginUserGeneral(req, ctx)
	if !errRsp.Status {
		ctx.JSON(errRsp.HttpStatus, server.errorResponse(errRsp.Error, errRsp.Message))
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}
