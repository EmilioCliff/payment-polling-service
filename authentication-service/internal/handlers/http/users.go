package http

import (
	"net/http"
	"time"

	"github.com/EmilioCliff/payment-polling-app/authentication-service/internal/repository"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/pkg"
	"github.com/gin-gonic/gin"
)

type registerUserRequest struct {
	FullName  string    `json:"full_name" binding:"required"`
	Email     string    `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	PaydUsername string `json:"payd_username" binding:"required"`
	PaydAccountID   string    `json:"payd_account_id" binding:"required"`
	PaydUsernameKey string    `json:"payd_username_key" binding:"required"`
	PaydPasswordKey string    `json:"payd_password_key" binding:"required"`
}

type registerUserResponse struct {
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *HTTPServer) handleRegisterUser(ctx *gin.Context) {
	var req registerUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rsp, err := s.UserRepository.CreateUser(ctx, repository.User{
		FullName: req.FullName,
		Email: req.Email,
		Password: req.Password,
		PaydUsername: req.PaydUsername,
		PaydAccountID: req.PaydAccountID,
		PaydUsernameKey: req.PaydUsernameKey,
		PaydPasswordKey: req.PaydPasswordKey,
	})
	if err != nil {
		ctx.JSON(s.convertPkgError(pkg.ErrorCode(err)), gin.H{"status_code": "Guesed this Value", "message": pkg.ErrorMessage(err)})
		return
	}

	ctx.JSON(http.StatusOK, registerUserResponse{
		FullName: rsp.FullName,
		Email: rsp.Email,
		CreatedAt: rsp.CreatedAt,
	})
}

type LoginUserRequest struct {
    Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginUserResponse struct {
    AccessToken  string    `json:"access_token"`
	ExpirationAt time.Time `json:"expiration_at"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"created_at"`
}

func (s *HTTPServer) handleLoginUser(ctx *gin.Context) {
	var req LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rsp, err := s.UserRepository.GetUser(ctx, req.Email)
	if err != nil {
		ctx.JSON(s.convertPkgError(pkg.ErrorCode(err)), gin.H{"status_code": "Guesed this Value", "message": pkg.ErrorMessage(err)})
		return
	}

	err = pkg.ComparePasswordAndHash(rsp.Password, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status_code": "Guesed this Value", "message": err})
		return
	}

	accessToken, err := s.maker.CreateToken(rsp.Email, s.config.TOKEN_DURATION)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status_code": "Guesed this Value", "message": err})
	}

	ctx.JSON(http.StatusOK, LoginUserResponse{
		AccessToken: accessToken,
		ExpirationAt: time.Now().Add(s.config.TOKEN_DURATION),
		FullName: rsp.FullName,
		Email: rsp.Email,
		CreatedAt: rsp.CreatedAt,
	})
}

func (s *HTTPServer) convertPkgError(err string) int {
    switch err {
    case pkg.ALREADY_EXISTS_ERROR:
        return http.StatusConflict
    case pkg.INTERNAL_ERROR:
        return http.StatusInternalServerError
    case pkg.INVALID_ERROR:
        return http.StatusBadRequest
    case pkg.NOT_FOUND_ERROR:
        return http.StatusNotFound
    case pkg.NOT_IMPLEMENTED_ERROR:
        return http.StatusNotImplemented
    case pkg.AUTHENTICATION_ERROR:
        return http.StatusUnauthorized
    default:
        return http.StatusInternalServerError
    }
}