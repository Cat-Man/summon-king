package middleware

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
)

const authTokenContextKey = "auth_token"

func InjectAuthToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		token := strings.TrimSpace(authHeader)
		if token != "" {
			parts := strings.Fields(token)
			if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
				token = parts[1]
			}
		}
		if token != "" {
			c.Set(authTokenContextKey, token)
			c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), authTokenContextKey, token))
		}
		c.Next()
	}
}

func GetAuthToken(c *gin.Context) string {
	if token, _ := c.Request.Context().Value(authTokenContextKey).(string); token != "" {
		return token
	}
	if value, ok := c.Get(authTokenContextKey); ok {
		if token, isString := value.(string); isString {
			return token
		}
	}
	return ""
}
