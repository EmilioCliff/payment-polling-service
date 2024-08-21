package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
)

func (server *Server) errorResponse(err error, msg string) gin.H {
	return gin.H{
		"message": msg,
		"error":   err.Error(),
	}
}

func (server *Server) grpcErrorResponse(code codes.Code, msg string) (int, gin.H) {
	switch code {
	case codes.Internal:
		return http.StatusInternalServerError, gin.H{
			"message": msg,
		}
	case codes.NotFound:
		return http.StatusNotFound, gin.H{
			"message": msg,
		}
	case codes.AlreadyExists:
		return http.StatusForbidden, gin.H{
			"message": msg,
		}
	case codes.Unauthenticated:
		return http.StatusUnauthorized, gin.H{
			"message": msg,
		}
	case codes.PermissionDenied:
		return http.StatusForbidden, gin.H{
			"message": msg,
		}
	default:
		return http.StatusInternalServerError, gin.H{
			"message": msg,
		}
	}
}
