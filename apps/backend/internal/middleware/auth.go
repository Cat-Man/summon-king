package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

const AccessTokenKey = "access_token"

func AuthToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := parseBearerToken(c.GetHeader("Authorization"))
		if token != "" {
			c.Set(AccessTokenKey, token)
		}
		c.Next()
	}
}

func GetAccessToken(c *gin.Context) string {
	if value, ok := c.Get(AccessTokenKey); ok {
		if token, ok := value.(string); ok {
			return token
		}
	}
	return ""
}

func parseBearerToken(header string) string {
	if len(header) < len("Bearer ") {
		return ""
	}
	if !strings.EqualFold(header[:len("Bearer ")], "Bearer ") {
		return ""
	}
	return strings.TrimSpace(header[len("Bearer "):])
}
