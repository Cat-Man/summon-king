package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
)

const idempotencyKeyContextKey = "idempotency_key"

func InjectIdempotencyKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-Idempotency-Key")
		if key != "" {
			c.Set(idempotencyKeyContextKey, key)
			c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), idempotencyKeyContextKey, key))
		}
		c.Next()
	}
}

func GetIdempotencyKey(c *gin.Context) string {
	if value, ok := c.Request.Context().Value(idempotencyKeyContextKey).(string); ok && value != "" {
		return value
	}
	if value, ok := c.Get(idempotencyKeyContextKey); ok {
		if key, isString := value.(string); isString {
			return key
		}
	}
	return ""
}
