package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	traceContextKey = "trace_id"
	traceIDHeader   = "X-Trace-ID"
)

func InjectTraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader(traceIDHeader)
		if traceID == "" {
			traceID = newTraceID()
		}

		c.Set(traceContextKey, traceID)
		c.Writer.Header().Set(traceIDHeader, traceID)
		c.Next()
	}
}

func GetTraceID(c *gin.Context) string {
	if value, ok := c.Get(traceContextKey); ok {
		if traceID, isString := value.(string); isString {
			return traceID
		}
	}
	return ""
}

func newTraceID() string {
	var bytes [16]byte
	if _, err := rand.Read(bytes[:]); err == nil {
		return hex.EncodeToString(bytes[:])
	}
	return hex.EncodeToString([]byte(time.Now().Format("20060102150405.000000000")))
}
