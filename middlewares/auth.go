package middlewares

import (
	"demo_webapp/controller"
	"demo_webapp/pkg/jwt"
	"github.com/gin-gonic/gin"
	"strings"
)

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			controller.ResponseError(c, controller.CodeNeedLogin)
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			controller.ResponseError(c, controller.CodeInvalidToken)
			c.Abort()
			return
		}
		cm, err := jwt.ParseToken(parts[1])
		if err != nil {
			controller.ResponseError(c, controller.CodeInvalidToken)
			c.Abort()
			return
		}

		c.Set(controller.CtxUserIDKey, cm.UserID)
		c.Next()
	}
}
