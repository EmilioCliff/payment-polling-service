package api

// CbCeJA4lYNxVMxb1tOC1 -- username
// ZNhXNJibSnkaW0eo3MzF2f1oFXfRnWJU1Z0JwfER -- password
import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/EmilioCliff/payment-polling-app/authentication-service/db/sqlc"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/utils"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type registerUserRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type registerResponse struct {
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type registerUserResponse struct {
	Status bool             `json:"status"`
	Data   registerResponse `json:"data"`
}

// registerUser handles the registration of a new user.
//
// It expects a JSON request body with the following fields:
// - fullName: the full name of the user (required)
// - email: the email of the user (required)
// - password: the password of the user (required)
//
// It returns a JSON response with the following fields:
// - status: a boolean indicating whether the registration was successful
// - data: a userResponse struct containing the user's full name, email, and creation timestamp
//
// If the request body is invalid, it returns a JSON response with a status code of 400 and an error message.
// If there is an error hashing the password, it returns a JSON response with a status code of 500 and an error message.
// If there is an error creating the user, it returns a JSON response with a status code of 500 and an error message.
// If the user already exists, it returns a JSON response with a status code of 403 and an error message.
func (server *Server) registerUser(ctx *gin.Context) {
	var req registerUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, server.errorResponse(err, "bad request to authApp"))
		return
	}

	hashPassword, err := utils.GenerateHashPassword(req.Password, server.config.HASH_COST)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, server.errorResponse(err, "error hashing password"))
		return
	}

	user, err := server.store.RegisterUser(ctx, db.RegisterUserParams{
		FullName: req.FullName,
		Email:    req.Email,
		Password: hashPassword,
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, server.errorResponse(err, "user already exists"))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, server.errorResponse(err, "error creating user"))
		return
	}

	rsp := registerUserResponse{
		Status: true,
		Data: registerResponse{
			FullName:  user.FullName,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	}

	ctx.JSON(http.StatusOK, rsp)
}

type loginUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	AccessToken  string    `json:"access_token"`
	ExpirationAt time.Time `json:"expiration_at"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"created_at"`
}

type loginUserResponse struct {
	Status bool          `json:"status"`
	Data   loginResponse `json:"data"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, server.errorResponse(err, "bad request to authApp"))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, server.errorResponse(err, "user not found"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, server.errorResponse(err, "error getting user by email"))
		return
	}

	err = utils.ComparePasswordAndHash(user.Password, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, server.errorResponse(err, "invalid credentials"))
		return
	}

	accessToken, err := server.maker.CreateToken(user.Email, server.config.TOKEN_DURATION)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, server.errorResponse(err, "error creating access token"))
		return
	}

	rsp := loginUserResponse{
		Status: true,
		Data: loginResponse{
			AccessToken:  accessToken,
			ExpirationAt: time.Now().Add(server.config.TOKEN_DURATION),
			FullName:     user.FullName,
			Email:        user.Email,
			CreatedAt:    user.CreatedAt,
		},
	}

	ctx.JSON(http.StatusOK, rsp)
}
