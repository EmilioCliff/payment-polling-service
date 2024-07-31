package api

import "github.com/gin-gonic/gin"

func (server *Server) errorResponse(err error, msg string) gin.H {
	return gin.H{
		"message": msg,
		"error":   err.Error(),
	}
}
