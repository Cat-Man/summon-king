package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

const authTokenContextKey = "auth_token"

func InjectAuthToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if token, ok := strings.CutPrefix(authHeader, "Bearer "); ok && token != "" {
			c.Set(authTokenContextKey, token)
		}
		c.Next()
	}
}

func GetAuthToken(c *gin.Context) string {
	if value, ok := c.Get(authTokenContextKey); ok {
		if token, isString := value.(string); isString {
			return token
		}
	}
	return ""
}
