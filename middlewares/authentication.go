package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/SirJager/tables/service/core/tokens"
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func BasicAuth(tokenBuilder tokens.Builder) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) < 1 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusForbidden, ErrorResponse{Error: err.Error()})
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusForbidden, ErrorResponse{Error: err.Error()})
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type: %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusForbidden, ErrorResponse{Error: err.Error()})
			return
		}

		access_token := fields[1]
		payload, err := tokenBuilder.VerifyToken(access_token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, ErrorResponse{Error: err.Error()})
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
