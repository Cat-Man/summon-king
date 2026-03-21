package middleware

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

const TraceIDKey = "trace_id"

func TraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = newTraceID()
		}

		c.Set(TraceIDKey, traceID)
		c.Writer.Header().Set("X-Trace-ID", traceID)
		c.Next()
	}
}

func GetTraceID(c *gin.Context) string {
	if value, ok := c.Get(TraceIDKey); ok {
		if traceID, ok := value.(string); ok {
			return traceID
		}
	}
	return ""
}

func newTraceID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "trace-fallback"
	}
	return hex.EncodeToString(buf)
}
