package services

import (
	"github.com/gin-gonic/gin"
)

type HttpInterface interface {
	RegisterUserViaHttp(RegisterUserRequest) (int, gin.H)
	LoginUserViaHttp(LoginUserRequest) (int, gin.H)
}
