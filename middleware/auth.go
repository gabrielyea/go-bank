package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gabriel/gabrielyea/go-bank/token"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey = "authorization"
	authType               = "bearer"
	authPayloadKey         = "auth_payload"
)

func authMiddleware(tMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header missing")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid auth header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		aType := strings.ToLower(fields[0])
		if aType != authType {
			err := errors.New("invalid auth type")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
		}

		accestoken := fields[1]
		payload, err := tMaker.VerifyToken(accestoken)
		if err != nil {
			err := errors.New("invalid token")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.Set(authPayloadKey, payload)
		ctx.Next()
	}
}
