package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func (server *Server) errorResponse(err error, msg string) gin.H {
	return gin.H{
		"status":  false,
		"message": fmt.Errorf(msg+": %w", err).Error(),
	}
}
