package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/EmilioCliff/payment-polling-app/authentication-service/utils"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationPayloadKey = "authorization_payload"
)

func authenticationMiddleware(tokenMaker utils.JWTMaker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizatonHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizatonHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "data": err.Error()})
			return
		}

		fields := strings.Fields(authorizatonHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "data": err.Error()})
			return
		}

		if fields[0] != "Bearer" {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "data": err.Error()})
			return
		}

		accessToken := fields[1]

		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "data": err.Error()})
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
