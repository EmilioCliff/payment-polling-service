package api

import "github.com/gin-gonic/gin"

func CORSMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token, Authorization, Accept")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET")
		ctx.Writer.Header().Set("Access-Control-Exposed-Headers", "Link")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}

		ctx.Next()
	}
}
