package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"go-alpha/response"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Failed("未登录", c)
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := response.ParseToken(tokenStr)
		if err != nil {
			response.Failed("登录已过期", c)
			c.Abort()
			return
		}

		c.Set("userId", claims.Id)
		c.Next()
	}
}
