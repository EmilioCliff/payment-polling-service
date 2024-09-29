package http

import (
	"net/http"
	"time"

	"github.com/EmilioCliff/payment-polling-app/authentication-service/internal/repository"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/pkg"
	"github.com/gin-gonic/gin"
)

type registerUserRequest struct {
	FullName        string `binding:"required" json:"full_name"`
	Email           string `binding:"required" json:"email"`
	Password        string `binding:"required" json:"password"`
	PaydUsername    string `binding:"required" json:"payd_username"`
	PaydAccountID   string `binding:"required" json:"payd_account_id"`
	PaydUsernameKey string `binding:"required" json:"payd_username_key"`
	PaydPasswordKey string `binding:"required" json:"payd_password_key"`
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
		FullName:        req.FullName,
		Email:           req.Email,
		Password:        req.Password,
		PaydUsername:    req.PaydUsername,
		PaydAccountID:   req.PaydAccountID,
		PaydUsernameKey: req.PaydUsernameKey,
		PaydPasswordKey: req.PaydPasswordKey,
	})
	if err != nil {
		statusCode := convertPkgError(pkg.ErrorCode(err))
		ctx.JSON(statusCode, gin.H{"status_code": statusCode, "message": pkg.ErrorMessage(err)})

		return
	}

	ctx.JSON(http.StatusOK, registerUserResponse{
		FullName:  rsp.FullName,
		Email:     rsp.Email,
		CreatedAt: rsp.CreatedAt,
	})
}

type LoginUserRequest struct {
	Email    string `binding:"required" json:"email"`
	Password string `binding:"required" json:"password"`
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
		statusCode := convertPkgError(pkg.ErrorCode(err))
		ctx.JSON(statusCode, gin.H{"status_code": statusCode, "message": pkg.ErrorMessage(err)})

		return
	}

	err = pkg.ComparePasswordAndHash(rsp.Password, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status_code": http.StatusUnauthorized, "message": err.Error()})

		return
	}

	accessToken, err := s.maker.CreateToken(rsp.Email, rsp.ID, s.config.TOKEN_DURATION)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status_code": http.StatusInternalServerError, "message": err.Error()})
	}

	ctx.JSON(http.StatusOK, LoginUserResponse{
		AccessToken:  accessToken,
		ExpirationAt: time.Now().Add(s.config.TOKEN_DURATION),
		FullName:     rsp.FullName,
		Email:        rsp.Email,
		CreatedAt:    rsp.CreatedAt,
	})
}

func convertPkgError(err string) int {
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
