package middlewares

import (
	"net/http"

	"github.com/SirJager/tables/config"
	"github.com/gin-gonic/gin"
)

const (
	AuthorizationRole           = "role"
	PUBLIC_ROLE                 = "PUBLIC"
	ADMIN_ROLE                  = "ADMIN"
	AuthorizationAdminSecretKey = "Admin-Secret"
)

func HasAdminPass() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// If the request as admin secret than it is admin and we will set isAdmin to true
		adminSecretKey := ctx.GetHeader(AuthorizationAdminSecretKey)
		c, _ := config.LoadConfig(".")
		if c.AdminSecretKey != adminSecretKey {
			ctx.AbortWithStatusJSON(http.StatusForbidden, ErrorResponse{Error: "you are not authorized to access this resource"})
			return
		}
		ctx.Next()
	}
}
