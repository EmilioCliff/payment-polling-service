package services

import "github.com/gin-gonic/gin"

type GrpcInterface interface {
	RegisterUserViagRPC(RegisterUserRequest) (int, gin.H)
	LoginUserViagRPC(LoginUserRequest) (int, gin.H)
}
