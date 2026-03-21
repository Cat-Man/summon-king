package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

const IdempotencyKey = "idempotency_key"

func IdempotencyKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := strings.TrimSpace(c.GetHeader("X-Idempotency-Key"))
		if key != "" {
			c.Set(IdempotencyKey, key)
		}
		c.Next()
	}
}

func GetIdempotencyKey(c *gin.Context) string {
	if value, ok := c.Get(IdempotencyKey); ok {
		if key, ok := value.(string); ok {
			return key
		}
	}
	return ""
}
